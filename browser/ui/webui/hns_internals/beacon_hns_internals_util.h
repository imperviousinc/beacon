#ifndef BEACON_BROWSER_UI_WEBUI_BEACON_HNS_INTERNALS_BEACON_HNS_INTERNALS_UTIL_H_
#define BEACON_BROWSER_UI_WEBUI_BEACON_HNS_INTERNALS_BEACON_HNS_INTERNALS_UTIL_H_

#include "base/callback.h"
#include "url/gurl.h"

class Browser;
class PrefService;

namespace beacon_hns_internals {
extern const char kBeaconHNSInternalsURL[];
extern const char kBeaconHNSInternalsURLShort[];

GURL GetServerURL(bool may_redirect);

}  // namespace beacon_hns_internals

#endif  // BEACON_BROWSER_UI_WEBUI_BEACON_HNS_INTERNALS_BEACON_HNS_INTERNALS_UTIL_H_
