// Modified for Beacon use
/* Copyright (c) 2019 The Brave Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#define CanAddURLToHistory CanAddURLToHistory_ChromiumImpl
#include "src/chrome/browser/history/history_utils.cc"
#undef CanAddURLToHistory

#include "beacon/common/url_constants.h"


bool CanAddURLToHistory(const GURL& url) {
  if (!CanAddURLToHistory_ChromiumImpl(url))
    return false;

  return !url.SchemeIs(content::kBeaconUIScheme);
}
