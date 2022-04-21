// Modified for Beacon use
/* Copyright (c) 2019 The Beacon Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#include "components/omnibox/browser/omnibox_client.h"
#include "components/omnibox/browser/omnibox_controller.h"
#include "url/gurl.h"

#if !BUILDFLAG(IS_IOS)
#include "content/public/common/url_constants.h"
#endif

// Keeping this it will be needed later
// although we're not overriding anything here yet.
class BeaconOmniboxController : public OmniboxController {
 public:
  BeaconOmniboxController(OmniboxEditModel* omnibox_edit_model,
                         OmniboxClient* client)
      : OmniboxController(omnibox_edit_model, client) {}
  BeaconOmniboxController(const BeaconOmniboxController&) = delete;
  BeaconOmniboxController& operator=(const BeaconOmniboxController&) = delete;
  ~BeaconOmniboxController() override = default;
};

namespace {
void BeaconAdjustTextForCopy(GURL* url) {
#if !BUILDFLAG(IS_IOS)
  if (url->scheme() == content::kChromeUIScheme) {
    GURL::Replacements replacements;
    replacements.SetSchemeStr(content::kBeaconUIScheme);
    *url = url->ReplaceComponents(replacements);
  }
#endif
}

}  // namespace

#define BEACON_ADJUST_TEXT_FOR_COPY \
  BeaconAdjustTextForCopy(url_from_text);

#define OmniboxController BeaconOmniboxController
#include "src/components/omnibox/browser/omnibox_edit_model.cc"
#undef OmniboxController
#undef BEACON_ADJUST_TEXT_FOR_COPY
