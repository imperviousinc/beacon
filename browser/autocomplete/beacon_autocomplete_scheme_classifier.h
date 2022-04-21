// Modified for Beacon use
/* Copyright (c) 2019 The Brave Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#ifndef BEACON_BROWSER_AUTOCOMPLETE_BEACON_AUTOCOMPLETE_SCHEME_CLASSIFIER_H_
#define BEACON_BROWSER_AUTOCOMPLETE_BEACON_AUTOCOMPLETE_SCHEME_CLASSIFIER_H_

#include <string>

#include "chrome/browser/autocomplete/chrome_autocomplete_scheme_classifier.h"

class BeaconAutocompleteSchemeClassifier
    : public ChromeAutocompleteSchemeClassifier {
 public:
  explicit BeaconAutocompleteSchemeClassifier(Profile* profile);
  BeaconAutocompleteSchemeClassifier(const BeaconAutocompleteSchemeClassifier&) =
      delete;
  BeaconAutocompleteSchemeClassifier& operator=(
      const BeaconAutocompleteSchemeClassifier&) = delete;
  ~BeaconAutocompleteSchemeClassifier() override;

  metrics::OmniboxInputType GetInputTypeForScheme(
      const std::string& scheme) const override;
};

#endif  // BEACON_BROWSER_AUTOCOMPLETE_BEACON_AUTOCOMPLETE_SCHEME_CLASSIFIER_H_

