#include <stdlib.h>
#include "hns.h"
#include "uv.h"
#include <assert.h>
#include "store.h"
#include "hsk.h"

#ifndef UNUSED
#define UNUSED(x) ((void)(x))
#endif

static void hns_uv_close_free(uv_handle_t *handle);
static hns_query *hns_queue_dequeue(hns_queue *queue);
static void hns_queue_uninit(hns_queue *queue);
static int hns_queue_init(hns_queue *queue);
static hns_queue *hns_queue_alloc();
static void hns_queue_free(hns_queue *queue);
static void hns_queue_enqueue(hns_queue *queue, hns_query *qry);

void hns_log(const char *fmt, ...) {
    printf("hns: ");

    va_list args;
    va_start(args, fmt);
    vprintf(fmt, args);
    va_end(args);
}

static int hsk_to_hns_err(int c) {
    switch (c) {
        case HSK_SUCCESS:
            return HNS_SUCCESS;
        case HSK_ETIMEOUT:
            return HNS_ETIMEOUT;
        case HSK_EBADARGS:
            return HNS_EBADARGS;
        case HSK_EFAILURE:
            return HNS_EFAILURE;
        case HSK_ENOMEM:
            return HNS_ENOMEM;
        default:
            return HNS_EUNKNOWN;
    }
}


static int hns_queue_init(hns_queue *queue) {
    if (uv_mutex_init(&queue->mutex))
        return HNS_EFAILURE;
    queue->head = NULL;
    queue->tail = NULL;

    return HNS_SUCCESS;
}

static void hns_queue_uninit(hns_queue *queue) {
    hns_query *current = queue->head;
    while (current) {
        hns_query *next = current->next;

        // clear request data in queue
        free(current->name);

        free(current);
        current = next;
    }

    uv_mutex_destroy(&queue->mutex);
}

static hns_queue *hns_queue_alloc() {
    hns_queue *queue = malloc(sizeof(hns_queue));
    if (!queue)
        return NULL;

    if (hns_queue_init(queue) != HNS_SUCCESS) {
        free(queue);
        return NULL;
    }

    return queue;
}

static void hns_queue_free(hns_queue *queue) {
    if (queue) {
        hns_queue_uninit(queue);
        free(queue);
    }
}

// Dequeue a query - thread-safe.
static hns_query *hns_queue_dequeue(hns_queue *queue) {
    uv_mutex_lock(&queue->mutex);

    hns_query *oldest = queue->head;
    if (oldest) {
        queue->head = oldest->next;
        oldest->next = NULL;
        if (queue->head)
            queue->head->prev = NULL; // Removed the prior request
        else {
            assert(queue->tail == oldest); // There was only one request
            queue->tail = NULL;
        }
    }

    uv_mutex_unlock(&queue->mutex);

    return oldest;
}

// Enqueue a query - thread-safe.
static void hns_queue_enqueue(hns_queue *queue, hns_query *qry) {
    uv_mutex_lock(&queue->mutex);

    if (!queue->tail) {
        // There were no requests queued; this one becomes head and tail
        assert(!queue->head);   // Invariant - set and cleared together
        queue->head = qry;
        queue->tail = qry;
    } else {
        // There are requests queued already, add this one to the tail
        queue->tail->next = qry;
        qry->prev = queue->tail;
        queue->tail = qry;
    }

    uv_mutex_unlock(&queue->mutex);
}

static uv_async_t * alloc_async(hns_ctx *ctx, uv_async_cb callback) {
    uv_async_t *async = malloc(sizeof(uv_async_t));
    if (!async) {
        return NULL;
    }
    async->data = NULL;

    // Initialize the async
    if (uv_async_init(ctx->loop, async, callback)) {
        free(async);
        return NULL;
    }

    async->data = (void *) ctx;
    return async;
}

static void free_async(uv_async_t *async) {
    if (async) {
        async->data = NULL;
        hns_uv_close_free((uv_handle_t *) async);
    }
}

static void free_timer(uv_timer_t *timer) {
    if (timer) {
        uv_timer_stop(timer);
        timer->data = NULL;
        hns_uv_close_free((uv_handle_t *) timer);
    }
}

static void after_close_free(uv_handle_t *handle) {
    free(handle);
}

static void hns_uv_close_free(uv_handle_t *handle) {
    if (handle)
        uv_close(handle, after_close_free);
}

static void hns_ctx_close_handles(hns_ctx *ctx) {
    free_async(ctx->exit_signal);
    free_async(ctx->queue_signal);
    free_timer(ctx->sync_timer);

    if (ctx->pool) {
        hsk_pool_close(ctx->pool);
        hsk_pool_free(ctx->pool);
    }

    ctx->exit_signal = NULL;
    ctx->queue_signal = NULL;
    ctx->sync_timer = NULL;
    ctx->pool = NULL;
}

void hns_ctx_destroy(hns_ctx *ctx) {
    hns_ctx_close_handles(ctx);

    if (ctx->pool_state) {
        uv_rwlock_destroy(&ctx->pool_state->lock);
        free(ctx->pool_state);
        ctx->pool_state = NULL;
    }

    if (ctx->loop) {
        while (uv_loop_close(ctx->loop) != 0) {
            uv_run(ctx->loop, UV_RUN_ONCE);
        }
        free(ctx->loop);
        ctx->loop = NULL;
    }

    if (ctx->queue) {
        hns_queue_free(ctx->queue);
        ctx->queue = NULL;
    }

    if (ctx->headers_file) {
        free(ctx->headers_file);
        ctx->headers_file = NULL;
    }

    free(ctx);
}

static void hns_call_cgo_cleanup(uv_work_t *req, int status) {
    UNUSED(status);

    hns_cgo_baton *baton = (hns_cgo_baton *) req->data;
    if (!baton)
        return;

    free(baton->name);
    if (baton->data_len > 0) {
        free(baton->data);
    }

    free(baton);
}

static void hns_call_cgo(uv_work_t *req) {
    hns_cgo_baton *baton = (hns_cgo_baton *) req->data;
    if (!baton)
        return;

    cgoAfterResolve(baton->name, baton->status, baton->exists, baton->data, baton->data_len, baton->ctx);
}

static void after_resolve(
        const char *name,
        int status,
        bool exists,
        const uint8_t *data,
        size_t data_len,
        const void *arg
) {
    hns_ctx *ctx = (hns_ctx *) arg;
    if (!ctx)
        return;

    hns_cgo_baton *baton = (hns_cgo_baton *) malloc(sizeof(hns_cgo_baton));
    baton->req.data = (void *) baton;
    baton->name = strdup(name);
    baton->status = hsk_to_hns_err(status);
    baton->exists = exists;
    baton->data_len = data_len;
    baton->ctx = ctx;

    if (data_len > 0) {
        baton->data = (uint8_t *) malloc(data_len);
        if (baton->data) {
            memcpy(baton->data, data, data_len);
        } else {
            baton->data_len = -1;
            baton->status = HNS_ENOMEM;
        }
    }

    uv_queue_work(ctx->loop, &baton->req, hns_call_cgo, hns_call_cgo_cleanup);
}

static bool has_active_peers(hns_ctx *ctx) {
    hsk_peer_t *peerIter, *next;
    for (peerIter = ctx->pool->head; peerIter; peerIter = next) {
        next = peerIter->next;
        if (peerIter->state == HSK_STATE_HANDSHAKE) {
            return true;
        }
    }

    return false;
}

static bool chain_ready(hns_ctx *ctx) {
    int64_t now = hsk_timedata_now(ctx->pool->chain.td);
    if (((int64_t) ctx->pool->chain.tip->time) < now - 21600)
        return false;

    return true;
}

static void resolve_name(hns_ctx *ctx, const char *name) {
    int rc = HNS_SUCCESS;

    if (!ctx->pool->chain.synced || !chain_ready(ctx)) {
        rc = HNS_ENOTSYNCED;
    } else if (!has_active_peers(ctx)) {
        rc = HNS_ENOPEERS;
    }

    if (rc == HNS_SUCCESS) {
        rc = hsk_pool_resolve(ctx->pool, name, after_resolve, (void *) ctx);
        if (rc == HSK_SUCCESS)
            return;

        rc = hsk_to_hns_err(rc);
    }

    hns_cgo_baton *baton = (hns_cgo_baton *) malloc(sizeof(hns_cgo_baton));
    baton->req.data = (void *) baton;
    baton->name = strdup(name);
    baton->status = rc;
    baton->exists = 0;
    baton->data_len = 0;
    baton->data = NULL;
    baton->ctx = ctx;

    uv_queue_work(ctx->loop, &baton->req, hns_call_cgo, hns_call_cgo_cleanup);
}

static void on_queue_signal(uv_async_t *async) {
    hns_ctx *ctx = (hns_ctx *) async->data;

    // Since uv_close() is async, it might be possible to process this event after
    // the ctx is destroyed but before the async is closed.
    if (!ctx)
        return;

    // Dequeue and process all events in the queue - libuv coalesces calls to
    // uv_async_send().
    hns_query *qry = hns_queue_dequeue(ctx->queue);
    while (qry) {
        // cgo callback
        hns_log("queue is processing name: %s\n", qry->name);
        resolve_name(ctx, qry->name);
        free(qry->name);
        free(qry);
        qry = hns_queue_dequeue(ctx->queue);
    }
}

static void on_exit_signal(uv_async_t *async) {
    hns_ctx *ctx = (hns_ctx *) async->data;

    // Should never get this after ctx is destroyed, the ctx can't be
    // destroyed until _close() completes.
    assert(ctx);
    hns_log("shutting down\n");
    hns_ctx_close_handles(ctx);
}

void hns_ctx_shutdown(hns_ctx *ctx) {
    uv_async_send(ctx->exit_signal);
}

void hns_resolve(hns_ctx *ctx, const char *name) {
    hns_query *q = (hns_query *) malloc(sizeof(hns_query));
    q->prev = NULL;
    q->next = NULL;
    q->ctx = ctx;
    q->name = strdup(name);

    hns_queue_enqueue(ctx->queue, q);
    uv_async_send(ctx->queue_signal);
}

void hns_ctx_set_id(hns_ctx *ctx, uint64_t id) {
    assert(ctx);
    ctx->id = id;
}

uint64_t hns_ctx_get_id(hns_ctx *ctx) {
    assert(ctx);
    return ctx->id;
}

static float chain_progress(hns_ctx *ctx) {
    double start = (double) ctx->pool->chain.genesis->time;
    double current = (double) ctx->pool->chain.tip->time - start;
    int64_t now = hsk_timedata_now(ctx->pool->chain.td);

    double end = (double) now - start - 40 * 60;
    return (float) (current/ end);
}

static void update_pool_state(hns_ctx *ctx) {
    float progress = chain_progress(ctx);
    bool ready = chain_ready(ctx);

    int total_peers = ctx->pool->size;
    int active_peers = 0;

    hsk_peer_t *peerIter, *next;
    for (peerIter = ctx->pool->head; peerIter; peerIter = next) {
        next = peerIter->next;
        if (peerIter->state == HSK_STATE_HANDSHAKE)
            active_peers++;
    }

    uv_rwlock_wrlock(&ctx->pool_state->lock);
    ctx->pool_state->chain_ready = ready;
    ctx->pool_state->chain_height = ctx->pool->chain.height;
    ctx->pool_state->sync_progress = progress;
    ctx->pool_state->total_peers = total_peers;
    ctx->pool_state->active_peers = active_peers;
    uv_rwlock_wrunlock(&ctx->pool_state->lock);

    // If there's enough proof-of-work
    // on top of the most recent root,
    // it should be safe to use it.
    uint32_t mod = (uint32_t)ctx->pool->chain.height % 36;
    if (mod >= 12) mod = 0;
    uint32_t root_height = (uint32_t)ctx->pool->chain.height - mod;
    const hsk_header_t *hdr = hsk_chain_get_by_height(&ctx->pool->chain, root_height);
    if (!hdr) {
        return;
    }

    const uint8_t *root = hdr->name_root;
    uv_rwlock_wrlock(&ctx->pool_state->lock);
    memcpy(&ctx->pool_state->name_root, root, 32);
    uv_rwlock_wrunlock(&ctx->pool_state->lock);
}

static void sync_timer_tick(uv_timer_t *handle) {
    hns_ctx *ctx = (hns_ctx *) handle->data;
    if (!ctx)
        return;

    update_pool_state(ctx);

    if (!ctx->headers_file)
        return;

    if (!hsk_chain_synced(&ctx->pool->chain)) {
        return;
    }

    uint32_t diff = (uint32_t) ctx->pool->chain.height - ctx->stored_height;
    if (ctx->stored_height != 0 && diff < 12) {
        return;
    }

    if (hns_write_chain(ctx, ctx->headers_file) == HNS_SUCCESS) {
        hns_log("block headers stored successfully\n");
        return;
    }

    hns_log("failed storing block headers\n");
}

hns_ctx *hns_ctx_create() {
    hns_ctx *ctx = (hns_ctx *) malloc(sizeof(hns_ctx));
    if (!ctx)
        return NULL;

    ctx->id = 0;
    ctx->queue = NULL;
    ctx->queue_signal = NULL;
    ctx->exit_signal = NULL;
    ctx->sync_timer = NULL;
    ctx->pool_state = NULL;
    ctx->stored_height = 0;
    ctx->headers_file = NULL;

    ctx->loop = (uv_loop_t *) malloc(sizeof(uv_loop_t));
    if (!ctx->loop)
        goto fail;

    if (uv_loop_init(ctx->loop) != 0)
        goto fail;

    ctx->queue = hns_queue_alloc();
    if (!ctx->queue) {
        goto fail;
    }

    ctx->queue_signal = alloc_async(ctx, on_queue_signal);
    if (!ctx->queue_signal)
        goto fail;

    ctx->exit_signal = alloc_async(ctx, on_exit_signal);
    if (!ctx->exit_signal)
        goto fail;

    ctx->sync_timer = (uv_timer_t *) malloc(sizeof(uv_timer_t));
    if (!ctx->sync_timer)
        goto fail;

    if (uv_timer_init(ctx->loop, ctx->sync_timer) != 0)
        goto fail;

    ctx->sync_timer->data = ctx;

    ctx->pool = hsk_pool_alloc(ctx->loop);
    if (!ctx->pool)
        goto fail;

    if (!hsk_pool_set_size(ctx->pool, 4))
        goto fail;

    if (!hsk_pool_set_agent(ctx->pool, "beacon"))
        goto fail;

    ctx->pool_state = (hns_pool_state *) malloc(sizeof(hns_pool_state));
    if (!ctx->pool_state)
        goto fail;

    ctx->pool_state->total_peers = 0;
    ctx->pool_state->active_peers = 0;
    ctx->pool_state->chain_height = 0;
    ctx->pool_state->sync_progress = 0;
    ctx->pool_state->chain_ready = false;
    memset(ctx->pool_state->name_root, 0, 32);
    assert(uv_rwlock_init(&ctx->pool_state->lock) == 0);

    return ctx;

    fail:
    hns_ctx_destroy(ctx);
    return NULL;
}

int hns_ctx_set_headers_file(hns_ctx *ctx, const char *fname) {
    ctx->headers_file = strdup(fname);
    return 0;
}

int hns_ctx_start(hns_ctx *ctx) {
    if (ctx->headers_file) {
        hns_read_chain(ctx,ctx->headers_file);
    }

    if (hsk_pool_open(ctx->pool) != HSK_SUCCESS) {
        hns_log("failed opening pool\n");
        return HNS_EFAILURE;
    }

    int rc = uv_timer_start(ctx->sync_timer, sync_timer_tick, 0, 500);
    if (rc != 0) {
        hns_log("failed starting timer: %s\n", uv_strerror(rc));
        return HNS_EFAILURE;
    }

    rc = uv_run(ctx->loop, UV_RUN_DEFAULT);
    if (rc != 0) {
        hns_log("uv run failed: %s\n", uv_strerror(rc));
        return HNS_EFAILURE;
    }

    return HNS_SUCCESS;
}

float hns_chain_progress(hns_ctx *ctx) {
    if (!ctx || !ctx->pool_state)
        return 0;

    uv_rwlock_rdlock(&ctx->pool_state->lock);
    float progress = ctx->pool_state->sync_progress;
    uv_rwlock_rdunlock(&ctx->pool_state->lock);
    return progress;
}

uint32_t hns_chain_height(hns_ctx *ctx) {
    if (!ctx || !ctx->pool_state)
        return 0;

    uv_rwlock_rdlock(&ctx->pool_state->lock);
    uint32_t height = ctx->pool_state->chain_height;
    uv_rwlock_rdunlock(&ctx->pool_state->lock);
    return height;
}

uint8_t* hns_chain_name_root(hns_ctx *ctx) {
    if (!ctx || !ctx->pool_state)
        return NULL;

    uint8_t* root = (uint8_t*)malloc(32);
    memset(root, 0, 32);
    uv_rwlock_rdlock(&ctx->pool_state->lock);
    memcpy(root, ctx->pool_state->name_root, 32);
    uv_rwlock_rdunlock(&ctx->pool_state->lock);
    return root;
}

bool hns_chain_ready(hns_ctx *ctx) {
    if (!ctx || !ctx->pool_state)
        return false;

    uv_rwlock_rdlock(&ctx->pool_state->lock);
    bool ready = ctx->pool_state->chain_ready;
    uv_rwlock_rdunlock(&ctx->pool_state->lock);
    return ready;
}

int hns_pool_total_peers(hns_ctx *ctx) {
    if (!ctx || !ctx->pool_state)
        return 0;

    uv_rwlock_rdlock(&ctx->pool_state->lock);
    int peers = ctx->pool_state->total_peers;
    uv_rwlock_rdunlock(&ctx->pool_state->lock);
    return peers;
}

int hns_pool_active_peers(hns_ctx *ctx) {
    if (!ctx || !ctx->pool_state) {
        return 0;
    }

    uv_rwlock_rdlock(&ctx->pool_state->lock);
    int active = ctx->pool_state->active_peers;
    uv_rwlock_rdunlock(&ctx->pool_state->lock);
    return active;
}
