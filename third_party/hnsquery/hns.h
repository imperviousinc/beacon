#ifndef HNSQ_HNS_H
#define HNSQ_HNS_H

#include <pool.h>
#include "uv.h"
#include <stdarg.h>

#define HNS_SUCCESS 0
#define HNS_ENOMEM 1
#define HNS_ETIMEOUT 2
#define HNS_EFAILURE 3
#define HNS_EBADARGS 4
#define HNS_ENOPEERS 5
#define HNS_ENOTSYNCED 6
#define HNS_EUNKNOWN 7

struct hns_query_s;
struct hns_pool_state_s;
struct hns_queue_s;

typedef struct hns_pool_state_s hns_pool_state;
typedef struct hns_query_s hns_query;
typedef struct hns_queue_s hns_queue;

typedef struct hns_ctx {
    // Unique id for this context
    uint64_t id;

    // Queue of requests received from cgo.
    hns_queue *queue;

    // Async used to signal libuv event loop to read
    // from queue
    uv_async_t *queue_signal;

    // Async used to signal libuv event loop to exit
    uv_async_t *exit_signal;

    // Timer to store block headers
    // and update pool state
    uv_timer_t *sync_timer;

    // Event loop
    uv_loop_t *loop;

    // Handshake pool
    hsk_pool_t *pool;

    // Thread safe pool state
    hns_pool_state *pool_state;

    // last stored height
    uint32_t stored_height;

    char *headers_file;
} hns_ctx;

typedef struct hns_cgo_baton {
    uv_work_t req;
    char * name;
    int status;
    bool exists;
    uint8_t * data;
    size_t data_len;
    hns_ctx *ctx;
} hns_cgo_baton;

extern void cgoAfterResolve(
        const char * name,
        int status,
        bool exists,
        const uint8_t * data,
        size_t data_len,
        const void * arg
);

// Thread-safe request queue
struct hns_queue_s {
    uv_mutex_t mutex;
    hns_query *head;
    hns_query *tail;
};

// Thread safe pool state
struct hns_pool_state_s {
    uv_rwlock_t lock;
    bool chain_ready;
    uint32_t chain_height;
    float sync_progress;
    int total_peers;
    int active_peers;
    uint8_t name_root[32];
};

// Request data from cgo
struct hns_query_s {
    hns_query *prev;
    hns_query *next;

    // The ctx that enqueued the request
    hns_ctx *ctx;
    char *name; // name to resolve
};

// Context create, start and destroy functions
// should be called from the same thread.
// hns_ctx_destroy can be called after
// the blocking hns_ctx_start returns
//
// To shutdown from a different thread use the thread-safe
// hns_ctx_shutdown function

// Creates a new context must free it with hns_ctx_destroy
hns_ctx *hns_ctx_create();

// Set file path to store block headers must be called before start
int hns_ctx_set_headers_file(hns_ctx *ctx, const char *fname);

void hns_ctx_set_id(hns_ctx *ctx, uint64_t id);

uint64_t hns_ctx_get_id(hns_ctx *ctx);

// Starts the context's event loop
int hns_ctx_start(hns_ctx *ctx);

// Frees the context memory
void hns_ctx_destroy(hns_ctx *ctx);

// Thread-safe - sends a shutdown signal to the context's even loop
void hns_ctx_shutdown(hns_ctx *ctx);

void hns_log(const char *fmt, ...);

// Thread safe - queues a name to be resolved
void hns_resolve(hns_ctx *ctx, const char *name);

// Thread safe - gets the chain sync progress
float hns_chain_progress(hns_ctx *ctx);

// Thread safe - get current block height
uint32_t hns_chain_height(hns_ctx *ctx);

// Thread safe - get current name root
uint8_t* hns_chain_name_root(hns_ctx *ctx);

// Thread safe - returns true if the chain is synced
// and the last block timestamp is within 6 hours
bool hns_chain_ready(hns_ctx *ctx);

// Thread safe - current total peers in the pool
int hns_pool_total_peers(hns_ctx *ctx);

// Thread safe - current active peers in the pool
int hns_pool_active_peers(hns_ctx *ctx);


#endif //HNSQ_HNS_H
