// Modified for Beacon use
/* Copyright 2019 The Brave Authors. All rights reserved.
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this file,
 * You can obtain one at http://mozilla.org/MPL/2.0/. */

#ifndef BEACON_BROWSER_UI_VIEWS_PROFILES_BEACON_PROFILE_MENU_VIEW_H_
#define BEACON_BROWSER_UI_VIEWS_PROFILES_BEACON_PROFILE_MENU_VIEW_H_

#include "chrome/browser/ui/views/profiles/profile_menu_view.h"

class BeaconProfileMenuView : public ProfileMenuView {
 public:
  BeaconProfileMenuView(const BeaconProfileMenuView&) = delete;
  BeaconProfileMenuView& operator=(const BeaconProfileMenuView&) = delete;

 private:
  friend class ProfileMenuView;

  using ProfileMenuView::ProfileMenuView;
  ~BeaconProfileMenuView() override = default;

  // Helper methods for building the menu.
  void BuildIdentity() override;
  void BuildAutofillButtons() override;
  void BuildSyncInfo() override;
  void BuildFeatureButtons() override;
  gfx::ImageSkia GetSyncIcon() const override;
};

#endif  // BEACON_BROWSER_UI_VIEWS_PROFILES_BEACON_PROFILE_MENU_VIEW_H_