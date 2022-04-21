#include "beacon/browser/ui/webui/beacon_web_ui_controller_factory.h"

#include <memory>

#include "base/feature_list.h"
#include "base/memory/ptr_util.h"
// #include "beacon/common/beacon_features.h"
// #include "beacon/common/pref_names.h"
// #include "beacon/common/webui_url_constants.h"
#include "build/build_config.h"
#include "chrome/browser/profiles/profile.h"
#include "chrome/common/url_constants.h"
#include "components/prefs/pref_service.h"
#include "content/public/browser/web_contents.h"
#include "url/gurl.h"
#include "base/logging.h"

#if !BUILDFLAG(IS_ANDROID)
#include "beacon/browser/ui/webui/beacon_welcome/beacon_welcome_ui.h"
#include "beacon/browser/ui/webui/hns_internals/beacon_hns_internals_ui.h"
#endif

using content::WebUI;
using content::WebUIController;

namespace {

// A function for creating a new WebUI. The caller owns the return value, which
// may be NULL (for example, if the URL refers to an non-existent extension).
typedef WebUIController* (*WebUIFactoryFunction)(WebUI* web_ui,
                                                 const GURL& url);

WebUIController* NewWebUI(WebUI* web_ui, const GURL& url) {
  auto host = url.host_piece();
//   Profile* profile = Profile::FromBrowserContext(
//       web_ui->GetWebContents()->GetBrowserContext());
  if (host == "welcome") {
    LOG(INFO) << "Hitting welcome ui page";
    return new BeaconWelcomeUI(web_ui);
  }

  if (host == "hns-internals") {
    return new BeaconHNSInternalsUI(web_ui);
  }
  
  return nullptr;
}

// Returns a function that can be used to create the right type of WebUI for a
// tab, based on its URL. Returns NULL if the URL doesn't have WebUI associated
// with it.
WebUIFactoryFunction GetWebUIFactoryFunction(WebUI* web_ui,
                                             const GURL& url) {
  if (url.host_piece() == "welcome" || url.host_piece() == "hns-internals") {
    return &NewWebUI;
  }

  return nullptr;
}

}  // namespace

WebUI::TypeID BeaconWebUIControllerFactory::GetWebUIType(
      content::BrowserContext* browser_context, const GURL& url) {
  WebUIFactoryFunction function = GetWebUIFactoryFunction(NULL, url);
  if (function) {
    return reinterpret_cast<WebUI::TypeID>(function);
  }
  return ChromeWebUIControllerFactory::GetWebUIType(browser_context, url);
}

std::unique_ptr<WebUIController>
BeaconWebUIControllerFactory::CreateWebUIControllerForURL(WebUI* web_ui,
                                                         const GURL& url) {
  WebUIFactoryFunction function = GetWebUIFactoryFunction(web_ui, url);
  if (!function) {
    return ChromeWebUIControllerFactory::CreateWebUIControllerForURL(
        web_ui, url);
  }

  return base::WrapUnique((*function)(web_ui, url));
}


// static
BeaconWebUIControllerFactory* BeaconWebUIControllerFactory::GetInstance() {
  return base::Singleton<BeaconWebUIControllerFactory>::get();
}

BeaconWebUIControllerFactory::BeaconWebUIControllerFactory() {
}

BeaconWebUIControllerFactory::~BeaconWebUIControllerFactory() {
}
