// Modified for Beacon use
/* Copyright (c) 2019 The Brave Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#include "beacon/browser/ui/omnibox/beacon_omnibox_client_impl.h"

#include <algorithm>

#include "base/cxx17_backports.h"
#include "base/metrics/histogram_macros.h"
#include "base/values.h"
#include "beacon/browser/autocomplete/beacon_autocomplete_scheme_classifier.h"
#include "chrome/browser/profiles/profile.h"
#include "chrome/browser/ui/omnibox/chrome_omnibox_client.h"
#include "chrome/browser/ui/omnibox/chrome_omnibox_edit_controller.h"
#include "components/omnibox/browser/autocomplete_match.h"
#include "components/prefs/pref_registry_simple.h"
#include "components/prefs/pref_service.h"


BeaconOmniboxClientImpl::BeaconOmniboxClientImpl(
    OmniboxEditController* controller,
    Profile* profile)
    : ChromeOmniboxClient(controller, profile),
      profile_(profile),
      scheme_classifier_(profile) { }

BeaconOmniboxClientImpl::~BeaconOmniboxClientImpl() {}

const AutocompleteSchemeClassifier&
BeaconOmniboxClientImpl::GetSchemeClassifier() const {
  return scheme_classifier_;
}
