#ifndef BEACON_NET_DNS_BEACON_RESOLVE_CONTEXT_H_
#define BEACON_NET_DNS_BEACON_RESOLVE_CONTEXT_H_

#include "net/base/net_export.h"
#include "net/dns/resolve_context.h"

namespace net {

class DnsSession;
class URLRequestContext;

class NET_EXPORT_PRIVATE BeaconResolveContext : public ResolveContext {
 public:
  BeaconResolveContext(URLRequestContext* url_request_context,
                      bool enable_caching);

  BeaconResolveContext(const BeaconResolveContext&) = delete;
  BeaconResolveContext& operator=(const BeaconResolveContext&) = delete;

  ~BeaconResolveContext() override;

  bool GetDohServerAvailability(size_t doh_server_index,
                                const DnsSession* session) const override;
  size_t NumAvailableDohServers(const DnsSession* session) const override;

 private:
  bool IsFirstProbeCompleted(const ServerStats& stat) const;
};

}  // namespace net

#endif  // BEACON_NET_DNS_BEACON_RESOLVE_CONTEXT_H_
