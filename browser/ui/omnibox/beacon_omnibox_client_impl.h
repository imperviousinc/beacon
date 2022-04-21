// Modified for Beacon use
/* Copyright (c) 2019 The Brave Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#ifndef BEACON_BROWSER_UI_OMNIBOX_BEACON_OMNIBOX_CLIENT_IMPL_H_
#define BEACON_BROWSER_UI_OMNIBOX_BEACON_OMNIBOX_CLIENT_IMPL_H_

#include "base/memory/raw_ptr.h"
#include "beacon/browser/autocomplete/beacon_autocomplete_scheme_classifier.h"
#include "chrome/browser/ui/omnibox/chrome_omnibox_client.h"

class OmniboxEditController;
class PrefRegistrySimple;
class Profile;

class BeaconOmniboxClientImpl : public ChromeOmniboxClient {
 public:
  BeaconOmniboxClientImpl(OmniboxEditController* controller, Profile* profile);
  BeaconOmniboxClientImpl(const BeaconOmniboxClientImpl&) = delete;
  BeaconOmniboxClientImpl& operator=(const BeaconOmniboxClientImpl&) = delete;
  ~BeaconOmniboxClientImpl() override;

  const AutocompleteSchemeClassifier& GetSchemeClassifier() const override;

 private:
  raw_ptr<Profile> profile_ = nullptr;
  BeaconAutocompleteSchemeClassifier scheme_classifier_;
};

#endif  // BEACON_BROWSER_UI_OMNIBOX_BEACON_OMNIBOX_CLIENT_IMPL_H_
