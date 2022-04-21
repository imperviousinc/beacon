// Modified for Beacon use
/* Copyright (c) 2020 The Brave Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#define BEACON_WEB_UI_IMPL AddRequestableScheme(kBeaconUIScheme);
#include "src/content/browser/webui/web_ui_impl.cc"
#undef BEACON_WEB_UI_IMPL
