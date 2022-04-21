#ifndef BEACON_BROWSER_UI_WEBUI_BEACON_WELCOME_HANDLER_H_
#define BEACON_BROWSER_UI_WEBUI_BEACON_WELCOME_HANDLER_H_

#include "content/public/browser/web_ui_message_handler.h"

namespace base {
class ListValue;
}

// Page handler for chrome://welcome.
class BeaconWelcomeHandler : public content::WebUIMessageHandler {
 public:
  BeaconWelcomeHandler();
  ~BeaconWelcomeHandler() override;
  BeaconWelcomeHandler(const BeaconWelcomeHandler&) = delete;
  BeaconWelcomeHandler& operator=(const BeaconWelcomeHandler&) = delete;

 private:
  void HandleInitialize(const base::ListValue* args);

  // content::WebUIMessageHandler:
  void RegisterMessages() override;
};

#endif  // BEACON_BROWSER_UI_WEBUI_BEACON_WELCOME_HANDLER_H_
