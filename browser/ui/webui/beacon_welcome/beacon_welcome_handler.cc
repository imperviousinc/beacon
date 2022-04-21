#include "beacon/browser/ui/webui/beacon_welcome/beacon_welcome_handler.h"
#include "beacon/browser/ui/webui/beacon_welcome/beacon_welcome_util.h"


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

BeaconWelcomeHandler::BeaconWelcomeHandler() = default;

BeaconWelcomeHandler::~BeaconWelcomeHandler() = default;

void BeaconWelcomeHandler::RegisterMessages() {
  web_ui()->RegisterDeprecatedMessageCallback(
      "initialize", base::BindRepeating(&BeaconWelcomeHandler::HandleInitialize,
                                        base::Unretained(this)));
}

void BeaconWelcomeHandler::HandleInitialize(const base::ListValue* args) {
  const auto& list = args->GetList();
  CHECK_EQ(1U, list.size());
  const std::string& callback_id = list[0].GetString();

  AllowJavascript();
  ResolveJavascriptCallback(
      base::Value(callback_id),
      base::Value(beacon_welcome::GetServerURL(true).spec()));
}
