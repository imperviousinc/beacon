// Modified for Beacon use
/* Copyright 2019 The Brave Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#include "beacon/browser/autocomplete/beacon_autocomplete_scheme_classifier.h"

#include <string>

#include "base/strings/string_util.h"
#include "beacon/common/url_constants.h"
#include "chrome/browser/profiles/profile.h"

BeaconAutocompleteSchemeClassifier::BeaconAutocompleteSchemeClassifier(
    Profile* profile)
    : ChromeAutocompleteSchemeClassifier(profile) { }

BeaconAutocompleteSchemeClassifier::~BeaconAutocompleteSchemeClassifier() {
}

// Without this override, typing in beacon:// URLs will search Google
metrics::OmniboxInputType
BeaconAutocompleteSchemeClassifier::GetInputTypeForScheme(
    const std::string& scheme) const {
  if (scheme.empty()) {
    return metrics::OmniboxInputType::EMPTY;
  }
  if (base::IsStringASCII(scheme) &&
      base::LowerCaseEqualsASCII(scheme, kBeaconUIScheme)) {
    return metrics::OmniboxInputType::URL;
  }

  return ChromeAutocompleteSchemeClassifier::GetInputTypeForScheme(scheme);
}
