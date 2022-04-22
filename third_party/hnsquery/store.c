#include "store.h"

#include "hsk.h"
#include "uv.h"
#include "bio.h"

#define HNS_RAW_HDR_SIZE 272
#define HNS_HDR_SIZE (HNS_RAW_HDR_SIZE + 2)

size_t hns_header_write(const hsk_header_t *hdr, uint8_t **data) {
    size_t s = 0;

    // Preheader.
    s += write_u32(data, hdr->nonce);
    s += write_u64(data, hdr->time);
    s += write_bytes(data, hdr->prev_block, 32);
    s += write_bytes(data, hdr->name_root, 32);

    // Subheader.
    s += write_bytes(data, hdr->extra_nonce, 24);
    s += write_bytes(data, hdr->reserved_root, 32);
    s += write_bytes(data, hdr->witness_root, 32);
    s += write_bytes(data, hdr->merkle_root, 32);
    s += write_u32(data, hdr->version);
    s += write_u32(data, hdr->bits);

    // Mask.
    s += write_bytes(data, hdr->mask, 32);

    // Height & work.
    s += write_u32(data, hdr->height);
    s += write_bytes(data, hdr->work, 32);

    // Size checksum.
    s += write_u16(data, HNS_RAW_HDR_SIZE);

    return s;
}

bool hns_header_read(uint8_t **data, size_t *data_len, hsk_header_t *hdr) {
    // Preheader.
    if (!read_u32(data, data_len, &hdr->nonce))
        return false;

    if (!read_u64(data, data_len, &hdr->time))
        return false;

    if (!read_bytes(data, data_len, hdr->prev_block, 32))
        return false;

    if (!read_bytes(data, data_len, hdr->name_root, 32))
        return false;

    // Subheader.
    if (!read_bytes(data, data_len, hdr->extra_nonce, 24))
        return false;

    if (!read_bytes(data, data_len, hdr->reserved_root, 32))
        return false;

    if (!read_bytes(data, data_len, hdr->witness_root, 32))
        return false;

    if (!read_bytes(data, data_len, hdr->merkle_root, 32))
        return false;

    if (!read_u32(data, data_len, &hdr->version))
        return false;

    if (!read_u32(data, data_len, &hdr->bits))
        return false;

    // Mask.
    if (!read_bytes(data, data_len, hdr->mask, 32))
        return false;

    // Height & work.
    if (!read_u32(data, data_len, &hdr->height))
        return false;

    if (!read_bytes(data, data_len, hdr->work, 32))
        return false;

    // Size checksum.
    uint16_t size = 0;
    if (!read_u16(data, data_len, &size)) {
        hns_log("failed reading header checksum\n");
        return false;
    }

    if (size != HNS_RAW_HDR_SIZE) {
        hns_log("header size checksum didn't match\n");
        return false;
    }

    return true;
}

bool hns_header_decode(const uint8_t *data, size_t data_len, hsk_header_t *hdr) {
    return hns_header_read((uint8_t **) &data, &data_len, hdr);
}

size_t hns_header_encode(const hsk_header_t *hdr, uint8_t *data) {
    return hns_header_write(hdr, &data);
}

int hns_store_header(hns_ctx *ctx, FILE *fptr, uint32_t height) {
    hsk_header_t *hdr = hsk_chain_get_by_height(&ctx->pool->chain, height);
    assert(hdr);

    uint8_t data[HNS_HDR_SIZE];
    memset(data, 0, HNS_HDR_SIZE);
    assert(hns_header_encode(hdr, data) == HNS_HDR_SIZE);

    if (fwrite(data, HNS_HDR_SIZE, 1, fptr) == 1)
        return HNS_SUCCESS;

    return HNS_EFAILURE;
}


int hns_write_chain(hns_ctx *ctx, const char *chain) {
    uint32_t height = (uint32_t) ctx->pool->chain.height;

    // remove older data
    remove(chain);

    FILE *fptr = fopen(chain, "wb");
    if (fptr == NULL) {
        printf("cannot open file %s for writing\n", chain);
        return HNS_EFAILURE;
    }

    hns_log("writing block headers\n");
    // i = 1 skip genesis
    for (uint32_t i = 1; i < height; i++) {
        if (hns_store_header(ctx, fptr, i) != HNS_SUCCESS) {
            fclose(fptr);
            return HNS_EFAILURE;
        }
    }

    fclose(fptr);

    ctx->stored_height = (uint32_t) ctx->pool->chain.height;

    return HNS_SUCCESS;
}

static bool hns_chain_has_work(const hsk_chain_t *chain) {
    return memcmp(chain->tip->work, HSK_CHAINWORK, 32) >= 0;
}

static void hns_maybe_sync(hns_ctx *ctx) {
    if (ctx->pool->chain.synced) {
        return;
    }

    int64_t now = hsk_timedata_now(ctx->pool->chain.td);

    if (((int64_t) ctx->pool->chain.tip->time) < now - HSK_MAX_TIP_AGE)
        return;

    if (!hns_chain_has_work(&ctx->pool->chain))
        return;

    hns_log("chain is fully synced\n");
    ctx->pool->chain.synced = true;
}

int hns_read_chain(hns_ctx *ctx, const char *name) {
    FILE *fptr = fopen(name, "rb");
    if (fptr == NULL) {
        return HNS_EFAILURE;
    }

    uint8_t data[HNS_HDR_SIZE];
    memset(data, 0, HNS_HDR_SIZE);
    uint32_t last_height = 0;

    for (;;) {
        if (fread(data, HNS_HDR_SIZE, 1, fptr) != 1) {
            break;
        }

        hsk_header_t *hdr = hsk_header_alloc();
        hsk_header_init(hdr);
        if (!hns_header_decode(data, HNS_HDR_SIZE, hdr)) {
            free(hdr);
            break;
        }

        if (hdr->height - 1 != last_height) {
            hns_log("failed reading remaining block headers file likely corrupted");
            free(hdr);
            break;
        }

        last_height = hdr->height;

        const uint8_t *hash = hsk_header_cache(hdr);
        if (!hsk_map_set(&ctx->pool->chain.hashes, hash, (void *) hdr)) {
            fclose(fptr);
            free(hdr);
            return HNS_ENOMEM;
        }

        if (!hsk_map_set(&ctx->pool->chain.heights, &hdr->height, (void *) hdr)) {
            fclose(fptr);
            free(hdr);
            hsk_map_del(&ctx->pool->chain.hashes, hash);
            return HNS_ENOMEM;
        }

        ctx->pool->chain.height = hdr->height;
        ctx->pool->chain.tip = hdr;
    }

    hns_log("restored to chain height %lld\n", ctx->pool->chain.height);
    hns_maybe_sync(ctx);

    fclose(fptr);

    ctx->stored_height = ctx->pool->chain.height;
    return HNS_SUCCESS;
}
