// Modified for Beacon use
/* Copyright 2019 The Brave Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#ifndef BEACON_COMMON_BEACON_CONTENT_CLIENT_H_
#define BEACON_COMMON_BEACON_CONTENT_CLIENT_H_

#include <vector>

#include "chrome/common/chrome_content_client.h"

class BeaconContentClient : public ChromeContentClient {
 public:
  BeaconContentClient();
  ~BeaconContentClient() override;

 private:
  // ChromeContentClinet overrides:
  void AddAdditionalSchemes(Schemes* schemes) override;
};

#endif  // BEACON_COMMON_BEACON_CONTENT_CLIENT_H_
