#ifndef BEACON_CHROMIUM_SRC_NET_DNS_RESOLVE_CONTEXT_H_
#define BEACON_CHROMIUM_SRC_NET_DNS_RESOLVE_CONTEXT_H_

namespace net {
class BeaconResolveContext;
}  // namespace net

#define GetDohServerAvailability virtual GetDohServerAvailability
#define NumAvailableDohServers virtual NumAvailableDohServers
#define BEACON_RESOLVE_CONTEXT_H \
 private:                       \
  friend class BeaconResolveContext;

#include "src/net/dns/resolve_context.h"
#undef GetDohServerAvailability
#undef NumAvailableDohServers
#undef BEACON_RESOLVE_CONTEXT_H

#endif  // BEACON_CHROMIUM_SRC_NET_DNS_RESOLVE_CONTEXT_H_
