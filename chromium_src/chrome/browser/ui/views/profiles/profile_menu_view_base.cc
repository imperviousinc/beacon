// Modified for Beacon use
/* Copyright (c) 2019 The Brave Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#include "beacon/browser/ui/views/profiles/beacon_profile_menu_view.h"
#include "chrome/browser/ui/views/profiles/profile_menu_view.h"

#define ProfileMenuView BeaconProfileMenuView
#include "src/chrome/browser/ui/views/profiles/profile_menu_view_base.cc"
#undef ProfileMenuView
