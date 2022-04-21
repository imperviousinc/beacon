// Modified for Beacon use
/* Copyright (c) 2019 The Brave Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at https://mozilla.org/MPL/2.0/. */

#define REGISTER_BEACON_SCHEMES_DISPLAY_ISOLATED_AND_NO_JS                      \
  WebString beacon_scheme(WebString::FromASCII(kBeaconUIScheme));               \
  WebSecurityPolicy::RegisterURLSchemeAsDisplayIsolated(beacon_scheme);         \
  WebSecurityPolicy::RegisterURLSchemeAsNotAllowingJavascriptURLs(              \
      beacon_scheme);
  
#include "src/content/renderer/render_thread_impl.cc"
