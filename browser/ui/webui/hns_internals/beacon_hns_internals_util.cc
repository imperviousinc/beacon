#include "beacon/browser/ui/webui/hns_internals/beacon_hns_internals_util.h"
#include "base/bind.h"
#include "base/callback.h"
#include "base/check.h"
#include "base/command_line.h"
#include "base/feature_list.h"
#include "base/location.h"
#include "base/memory/raw_ptr.h"
#include "base/memory/scoped_refptr.h"
#include "base/metrics/histogram_functions.h"
#include "base/strings/string_number_conversions.h"
#include "base/strings/stringprintf.h"
#include "base/task/sequenced_task_runner.h"
#include "base/threading/sequenced_task_runner_handle.h"
#include "chrome/browser/browser_process.h"
#include "chrome/browser/net/system_network_context_manager.h"
#include "chrome/browser/profiles/profile.h"
#include "chrome/browser/ui/browser.h"
#include "chrome/browser/ui/browser_list.h"
#include "chrome/browser/ui/browser_list_observer.h"
#include "chrome/browser/ui/browser_tabstrip.h"
#include "chrome/browser/ui/ui_features.h"
#include "chrome/common/chrome_switches.h"
#include "chrome/common/chrome_version.h"
#include "chrome/common/pref_names.h"
#include "chrome/common/webui_url_constants.h"
#include "components/prefs/pref_service.h"
#include "content/public/browser/browser_task_traits.h"
#include "content/public/browser/browser_thread.h"
#include "net/base/url_util.h"
#include "net/http/http_util.h"
#include "services/network/public/cpp/resource_request.h"
#include "services/network/public/cpp/simple_url_loader.h"
#include "services/network/public/mojom/url_response_head.mojom.h"
#include "url/gurl.h"

namespace beacon_hns_internals {

const char kBeaconHNSInternalsURL[] = "http://127.0.0.1:44962/resources/hns-internals.html";
const char kBeaconHNSInternalsURLShort[] = "127.0.0.1";

GURL GetServerURL(bool may_redirect) {
  return may_redirect
             ? net::AppendQueryParameter(
                   GURL(kBeaconHNSInternalsURL), "version",
                   base::NumberToString(CHROME_VERSION_MAJOR))
             : GURL(kBeaconHNSInternalsURL)
                   .Resolve(base::StringPrintf("m%d", CHROME_VERSION_MAJOR));
}

}  // namespace beacon_hns_internals
