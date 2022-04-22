#ifndef CGO_BUILD

_Pragma("GCC diagnostic ignored \"-Wunused-parameter\"")

#include <stdio.h>
#include "hsk.h"
#include "hns.h"
#include <unistd.h>

void cgoAfterResolve(
        const char *name,
        int status,
        bool exists,
        const uint8_t *data,
        size_t data_len,
        const void *arg
) {
    printf("cgo: received name: %s\n", name);
}

void hns_thread(void *arg) {
    hns_ctx *ctx = (hns_ctx *) arg;
    assert(ctx);

    hns_ctx_start(ctx);
    hns_ctx_destroy(ctx);
}

void app_thread(void *arg) {
    hns_ctx *ctx = (hns_ctx *) arg;
    assert(ctx);

    const char *names[] = {"proofofconcept", "schematic", "3b"};
    const int len = 3;

    // attempt to resolve without waiting for ctx
    // to be ready
    for (int i = 0; i < len; i++) {
        hns_resolve(ctx, names[i]);
    }

    while (!hns_chain_ready(ctx)) {
        sleep(1);
    }

    // shutdown after 10 seconds
    sleep(10);
    hns_ctx_shutdown(ctx);
}

void app_thread2(void *arg) {
    hns_ctx *ctx = (hns_ctx *) arg;
    assert(ctx);

    const char *names[] = {"nb", "impervious", "forever"};
    const int len = 3;

    while (!hns_chain_ready(ctx)) {
        printf("height: %d, progress: %f, ready: %d, peers: (total: %d, active: %d)\n",
               hns_chain_height(ctx),
               hns_chain_progress(ctx),
               hns_chain_ready(ctx),
               hns_pool_total_peers(ctx),
               hns_pool_active_peers(ctx));

        sleep(1);
    }

    printf("height: %d, progress: %f, ready: %d, peers: (total: %d, active: %d)\n",
           hns_chain_height(ctx),
           hns_chain_progress(ctx),
           hns_chain_ready(ctx),
           hns_pool_total_peers(ctx),
           hns_pool_active_peers(ctx));

    for (int i = 0; i < len; i++) {
        hns_resolve(ctx, names[i]);
    }
    sleep(2);
}

int main() {
    hns_ctx *ctx = hns_ctx_create();
    assert(ctx);

    uv_thread_t hns_id;
    uv_thread_t app_id;
    uv_thread_t app_id2;

    uv_thread_create(&hns_id, hns_thread, ctx);
    uv_thread_create(&app_id, app_thread, ctx);
    uv_thread_create(&app_id2, app_thread2, ctx);

    uv_thread_join(&hns_id);
    uv_thread_join(&app_id);
    uv_thread_join(&app_id2);

    return 0;
}

#endif
