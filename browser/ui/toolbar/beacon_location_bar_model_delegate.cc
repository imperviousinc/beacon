// Modified for Beacon use
/* Copyright (c) 2019 The Brave Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#include "beacon/browser/ui/toolbar/beacon_location_bar_model_delegate.h"

#include "base/strings/utf_string_conversions.h"
#include "beacon/common/url_constants.h"
#include "chrome/browser/profiles/profile.h"
#include "chrome/browser/ui/browser.h"
#include "extensions/buildflags/buildflags.h"

BeaconLocationBarModelDelegate::BeaconLocationBarModelDelegate(Browser* browser)
    : BrowserLocationBarModelDelegate(browser) {}

BeaconLocationBarModelDelegate::~BeaconLocationBarModelDelegate() {}

// static
void BeaconLocationBarModelDelegate::FormattedStringFromURL(
    const GURL& url,
    std::u16string* new_formatted_url) {
  if (url.SchemeIs("chrome")) {
    base::ReplaceFirstSubstringAfterOffset(new_formatted_url, 0, u"chrome://",
                                           u"beacon://");
  }
}

std::u16string
BeaconLocationBarModelDelegate::FormattedStringWithEquivalentMeaning(
    const GURL& url,
    const std::u16string& formatted_url) const {
  std::u16string new_formatted_url =
      BrowserLocationBarModelDelegate::FormattedStringWithEquivalentMeaning(
          url, formatted_url);
  BeaconLocationBarModelDelegate::FormattedStringFromURL(url,
                                                        &new_formatted_url);
  return new_formatted_url;
}
