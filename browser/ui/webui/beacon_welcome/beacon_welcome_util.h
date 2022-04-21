// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#ifndef BEACON_BROWSER_UI_WEBUI_BEACON_WELCOME_BEACON_WELCOME_UTIL_H_
#define BEACON_BROWSER_UI_WEBUI_BEACON_WELCOME_BEACON_WELCOME_UTIL_H_

#include "base/callback.h"
#include "url/gurl.h"

class Browser;
class PrefService;

namespace beacon_welcome {
extern const char kBeaconWelcomeURL[];
extern const char kBeaconWelcomeURLShort[];

GURL GetServerURL(bool may_redirect);

}  // namespace beacon_welcome

#endif  // BEACON_BROWSER_UI_WEBUI_BEACON_WELCOME_BEACON_WELCOME_UTIL_H_
