/* Copyright (c) 2019 The Brave Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#ifndef BEACON_BROWSER_UI_WEBUI_BEACON_WEB_UI_CONTROLLER_FACTORY_H_
#define BEACON_BROWSER_UI_WEBUI_BEACON_WEB_UI_CONTROLLER_FACTORY_H_

#include <memory>

#include "chrome/browser/ui/webui/chrome_web_ui_controller_factory.h"

namespace base {
class RefCountedMemory;
}

class BeaconWebUIControllerFactory : public ChromeWebUIControllerFactory {
 public:
  BeaconWebUIControllerFactory(const BeaconWebUIControllerFactory&) = delete;
  BeaconWebUIControllerFactory& operator=(const BeaconWebUIControllerFactory&) =
      delete;

  content::WebUI::TypeID GetWebUIType(content::BrowserContext* browser_context,
                                      const GURL& url) override;
  std::unique_ptr<content::WebUIController> CreateWebUIControllerForURL(
      content::WebUI* web_ui,
      const GURL& url) override;

  static BeaconWebUIControllerFactory* GetInstance();

 protected:
  friend struct base::DefaultSingletonTraits<BeaconWebUIControllerFactory>;

  BeaconWebUIControllerFactory();
  ~BeaconWebUIControllerFactory() override;
};

#endif  // BEACON_BROWSER_UI_WEBUI_BEACON_WEB_UI_CONTROLLER_FACTORY_H_
