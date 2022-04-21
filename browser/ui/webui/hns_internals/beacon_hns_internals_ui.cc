#include "beacon/browser/ui/webui/hns_internals/beacon_hns_internals_ui.h"
#include "beacon/browser/ui/webui/hns_internals/beacon_hns_internals_util.h"
#include "beacon/browser/ui/webui/hns_internals/beacon_hns_internals_handler.h"

#include "base/feature_list.h"
#include "base/strings/stringprintf.h"
#include "base/version.h"
#include "chrome/browser/browser_features.h"
#include "chrome/browser/profiles/profile.h"
#include "chrome/browser/ui/hats/hats_service.h"
#include "chrome/browser/ui/hats/hats_service_factory.h"
#include "chrome/browser/ui/ui_features.h"
#include "chrome/browser/ui/webui/browser_command/browser_command_handler.h"
#include "chrome/browser/ui/webui/webui_util.h"
#include "chrome/common/chrome_features.h"
#include "chrome/common/pref_names.h"
#include "chrome/common/webui_url_constants.h"
#include "chrome/grit/generated_resources.h"
#include "chrome/grit/theme_resources.h"
#include "chrome/grit/beacon_hns_internals_resources.h"
#include "chrome/grit/beacon_hns_internals_resources_map.h"
#include "components/prefs/pref_registry_simple.h"
#include "components/prefs/pref_service.h"
#include "components/strings/grit/components_strings.h"
#include "content/public/browser/web_ui_data_source.h"
#include "services/network/public/mojom/content_security_policy.mojom.h"
#include "ui/base/resource/resource_bundle.h"
#include "ui/base/webui/web_ui_util.h"

namespace {

content::WebUIDataSource* CreateBeaconHNSInternalsUIHtmlSource(Profile* profile) {
  content::WebUIDataSource* source =
      content::WebUIDataSource::Create("hns-internals");

  webui::SetupWebUIDataSource(
      source, base::make_span(kBeaconHnsInternalsResources, kBeaconHnsInternalsResourcesSize),
      IDR_BEACON_HNS_INTERNALS_BEACON_HNS_INTERNALS_HTML);
 
  // Allow embedding of iframe from 127.0.0.1
  source->OverrideContentSecurityPolicy(
      network::mojom::CSPDirectiveName::ChildSrc,
      base::StringPrintf("child-src chrome://test http: %s;",
                         beacon_hns_internals::kBeaconHNSInternalsURLShort));
  return source;
}

}  // namespace

// static
void BeaconHNSInternalsUI::RegisterLocalStatePrefs(PrefRegistrySimple* registry) {
}

BeaconHNSInternalsUI::BeaconHNSInternalsUI(content::WebUI* web_ui)
    : ui::MojoWebUIController(web_ui, /*enable_chrome_send=*/true),
      browser_command_factory_receiver_(this),
      profile_(Profile::FromWebUI(web_ui)) {
  content::WebUIDataSource* source = CreateBeaconHNSInternalsUIHtmlSource(profile_);
  content::WebUIDataSource::Add(profile_, source);
  web_ui->AddMessageHandler(std::make_unique<BeaconHNSInternalsHandler>());
}

// static
base::RefCountedMemory* BeaconHNSInternalsUI::GetFaviconResourceBytes(
    ui::ResourceScaleFactor scale_factor) {
  return static_cast<base::RefCountedMemory*>(
      ui::ResourceBundle::GetSharedInstance().LoadDataResourceBytesForScale(
          IDR_NTP_FAVICON, scale_factor));
}

WEB_UI_CONTROLLER_TYPE_IMPL(BeaconHNSInternalsUI)

void BeaconHNSInternalsUI::BindInterface(
    mojo::PendingReceiver<browser_command::mojom::CommandHandlerFactory>
        pending_receiver) {
  if (browser_command_factory_receiver_.is_bound())
    browser_command_factory_receiver_.reset();
  browser_command_factory_receiver_.Bind(std::move(pending_receiver));
}

void BeaconHNSInternalsUI::CreateBrowserCommandHandler(
    mojo::PendingReceiver<browser_command::mojom::CommandHandler>
        pending_handler) {
}

BeaconHNSInternalsUI::~BeaconHNSInternalsUI() = default;
