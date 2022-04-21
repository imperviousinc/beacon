#include "beacon/browser/ui/webui/beacon_web_ui_controller_factory.h"
#include "build/chromeos_buildflags.h"

#define BEACON_CHROME_WEBUI_CONTROLLER_FACTORY \
  return BeaconWebUIControllerFactory::GetInstance();

#include "src/chrome/browser/ui/webui/chrome_web_ui_controller_factory.cc"
#undef BEACON_CHROME_WEBUI_CONTROLLER_FACTORY
