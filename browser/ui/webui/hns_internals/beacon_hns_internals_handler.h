#ifndef BEACON_BROWSER_UI_WEBUI_BEACON_HNS_INTERNALS_HANDLER_H_
#define BEACON_BROWSER_UI_WEBUI_BEACON_HNS_INTERNALS_HANDLER_H_

#include "content/public/browser/web_ui_message_handler.h"

namespace base {
class ListValue;
}

// Page handler for chrome://hns-internals.
class BeaconHNSInternalsHandler : public content::WebUIMessageHandler {
 public:
  BeaconHNSInternalsHandler();
  ~BeaconHNSInternalsHandler() override;
  BeaconHNSInternalsHandler(const BeaconHNSInternalsHandler&) = delete;
  BeaconHNSInternalsHandler& operator=(const BeaconHNSInternalsHandler&) = delete;

 private:
  void HandleInitialize(const base::ListValue* args);

  // content::WebUIMessageHandler:
  void RegisterMessages() override;
};

#endif  // BEACON_BROWSER_UI_WEBUI_BEACON_HNS_INTERNALS_HANDLER_H_
