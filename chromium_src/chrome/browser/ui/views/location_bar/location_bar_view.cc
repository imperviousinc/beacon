// Modified for Beacon use
/* Copyright (c) 2019 The Brave Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#include "beacon/browser/ui/omnibox/beacon_omnibox_client_impl.h"

#define ChromeOmniboxClient BeaconOmniboxClientImpl
#include "src/chrome/browser/ui/views/location_bar/location_bar_view.cc"
#undef ChromeOmniboxClient
