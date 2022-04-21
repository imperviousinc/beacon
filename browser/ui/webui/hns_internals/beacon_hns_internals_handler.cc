#include "beacon/browser/ui/webui/hns_internals/beacon_hns_internals_handler.h"
#include "beacon/browser/ui/webui/hns_internals/beacon_hns_internals_util.h"


#include "base/bind.h"
#include "base/metrics/histogram_functions.h"
#include "base/values.h"
#include "chrome/browser/browser_process.h"
#include "chrome/browser/profiles/profile.h"
#include "chrome/browser/ui/browser.h"
#include "chrome/browser/ui/browser_finder.h"
#include "chrome/browser/ui/browser_tabstrip.h"
#include "chrome/common/chrome_version.h"
#include "chrome/common/pref_names.h"
#include "chrome/common/webui_url_constants.h"
#include "components/bookmarks/common/bookmark_pref_names.h"
#include "components/prefs/pref_change_registrar.h"
#include "components/prefs/pref_service.h"
#include "url/gurl.h"

BeaconHNSInternalsHandler::BeaconHNSInternalsHandler() = default;

BeaconHNSInternalsHandler::~BeaconHNSInternalsHandler() = default;

void BeaconHNSInternalsHandler::RegisterMessages() {
  web_ui()->RegisterDeprecatedMessageCallback(
      "initialize", base::BindRepeating(&BeaconHNSInternalsHandler::HandleInitialize,
                                        base::Unretained(this)));
}

void BeaconHNSInternalsHandler::HandleInitialize(const base::ListValue* args) {
  const auto& list = args->GetList();
  CHECK_EQ(1U, list.size());
  const std::string& callback_id = list[0].GetString();

  AllowJavascript();
  ResolveJavascriptCallback(
      base::Value(callback_id),
      base::Value(beacon_hns_internals::GetServerURL(true).spec()));
}
