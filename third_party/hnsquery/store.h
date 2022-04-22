#ifndef HNSQ_STORE_H
#define HNSQ_STORE_H

#include "hns.h"

int hns_read_chain(hns_ctx *ctx, const char *filename);
int hns_write_chain(hns_ctx *ctx, const char *filename);

#endif //HNSQ_STORE_H
