// based on brave-core/net/dns
#include "base/strings/string_util.h"
#include "net/dns/dns_config.h"
#include "net/dns/dns_server_iterator.h"
#include "net/base/url_util.h"

// This allows dotless names to be resolvable with DoH
#define BEACON_DISABLE_SUFFIX_SEARCH qnames_.push_back(labeled_hostname); \
  return OK;

namespace {

bool GetNextIndex(const std::string& hostname,
                  const net::DnsConfig& config,
                  net::DnsServerIterator* dns_server_iterator,
                  size_t* doh_server_index) {
  base::StringPiece server =
      config.doh_config.servers()[*doh_server_index].server_template();

  // Brave uses this techniuqe for .eth but here we execlude
  // all ICANN TLDs. 
  // TODO: Non-unique hostnames relies on Chromium's public 
  // suffix list so we should define our own method
  // since Handshake introduces name collisions.
  while (server == "https://hs.dnssec.dev/dns-query" && 
    !net::IsHostnameNonUnique(hostname)) {
    // No next available index to attempt.
    if (!dns_server_iterator->AttemptAvailable()) {
      return false;
    }

    *doh_server_index = dns_server_iterator->GetNextAttemptIndex();
    server = config.doh_config.servers()[*doh_server_index].server_template();
  }

  return true;
}

}  // namespace

#define BEACON_MAKE_HTTP_ATTEMPT                                       \
  if (!GetNextIndex(hostname_, session_.get()->config(),              \
                    dns_server_iterator_.get(), &doh_server_index)) { \
    return AttemptResult(ERR_BLOCKED_BY_CLIENT, nullptr);             \
  }

#include "src/net/dns/dns_transaction.cc"
#undef BEACON_MAKE_HTTP_ATTEMPT
#undef BEACON_DISABLE_SUFFIX_SEARCH
