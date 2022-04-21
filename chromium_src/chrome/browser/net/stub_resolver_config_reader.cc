#include "components/prefs/pref_service.h"

#include "base/strings/strcat.h"
#include "net/dns/public/dns_over_https_config.h"

namespace {

void AddDoHServers(net::DnsOverHttpsConfig* doh_config,
                   PrefService* local_state,
                   bool force_check_parental_controls_for_automatic_mode) {
  if (force_check_parental_controls_for_automatic_mode)
    return;

  std::string doh_config_string = doh_config->ToString();

  if (doh_config_string.find("https://hs.dnssec.dev/dns-query") == std::string::npos) {
    doh_config_string =
        base::StrCat({"https://hs.dnssec.dev/dns-query", " ", doh_config_string});
  }

  *doh_config = net::DnsOverHttpsConfig::FromStringLax(doh_config_string);
}

}  // namespace

#define BEACON_GET_AND_UPDATE_CONFIGURATION \
  AddDoHServers(&doh_config, local_state_, \
                force_check_parental_controls_for_automatic_mode);

#include "src/chrome/browser/net/stub_resolver_config_reader.cc"
#undef BEACON_GET_AND_UPDATE_CONFIGURATION
