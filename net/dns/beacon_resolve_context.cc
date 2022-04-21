// based on brave-core/net/dns
#include "beacon/net/dns/beacon_resolve_context.h"

#include <string>

#include "net/dns/dns_session.h"

namespace {

bool IsHandshakeResolver(const std::string& server) {
  return server == "https://hs.dnssec.dev/dns-query";
}

}  // namespace

namespace net {

BeaconResolveContext::BeaconResolveContext(URLRequestContext* url_request_context,
                                         bool enable_caching)
    : ResolveContext(url_request_context, enable_caching) {}

BeaconResolveContext::~BeaconResolveContext() = default;

bool BeaconResolveContext::IsFirstProbeCompleted(const ServerStats& stat) const {
  return !(stat.last_failure_count == 0 &&
           stat.current_connection_success == false);
}

bool BeaconResolveContext::GetDohServerAvailability(
    size_t doh_server_index,
    const DnsSession* session) const {
  // TODO: This doesn't seem to be effective or some other cases also cause NX. 
  // Comment from Brave:
  // Return HNS resolver as available before the first probe is
  // completed. It is to avoid falling back to non-secure DNS servers before
  // the first probe is completed when users using automatic mode, which will
  // lead to an error page with HOSTNAME_NOT_RESOLVED error right after user
  // opt-in from the interstitial page.
  if (IsHandshakeResolver(session->config()
                                     .doh_config.servers()[doh_server_index]
                                     .server_template()) &&
      !IsFirstProbeCompleted(doh_server_stats_[doh_server_index]))
    return true;

  return ResolveContext::GetDohServerAvailability(doh_server_index, session);
}

size_t BeaconResolveContext::NumAvailableDohServers(
    const DnsSession* session) const {
  size_t num = 0;

  // Treat handshake resolver as available before the first probe is
  // completed.
  for (size_t i = 0; i < doh_server_stats_.size(); i++) {
    if (IsHandshakeResolver(
            session->config().doh_config.servers()[i].server_template()) &&
        !IsFirstProbeCompleted(doh_server_stats_[i]))
      num++;
  }

  return num + ResolveContext::NumAvailableDohServers(session);
}

}  // namespace net
