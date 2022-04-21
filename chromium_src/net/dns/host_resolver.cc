#include "beacon/net/dns/beacon_resolve_context.h"
#include "net/dns/context_host_resolver.h"

#define ResolveContext BeaconResolveContext
#include "src/net/dns/host_resolver.cc"
#undef ResolveContext
