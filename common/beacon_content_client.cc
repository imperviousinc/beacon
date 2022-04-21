// Modified for Beacon use
/* Copyright 2019 The Brave Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#include "beacon/common/beacon_content_client.h"

#include <string>

#include "base/memory/ref_counted_memory.h"
#include "components/grit/components_resources.h"
#include "content/public/common/url_constants.h"
#include "ui/base/resource/resource_bundle.h"


BeaconContentClient::BeaconContentClient() {}

BeaconContentClient::~BeaconContentClient() {}

void BeaconContentClient::AddAdditionalSchemes(Schemes* schemes) {
  ChromeContentClient::AddAdditionalSchemes(schemes);
  schemes->standard_schemes.push_back(content::kBeaconUIScheme);
  schemes->secure_schemes.push_back(content::kBeaconUIScheme);
  schemes->cors_enabled_schemes.push_back(content::kBeaconUIScheme);
  schemes->savable_schemes.push_back(content::kBeaconUIScheme);
}
