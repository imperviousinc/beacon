#ifndef BEACON_CHROMIUM_SRC_COMPONENTS_SECURITY_INTERSTITIALS_CORE_BASE_SAFE_BROWSING_ERROR_UI_H_
#define BEACON_CHROMIUM_SRC_COMPONENTS_SECURITY_INTERSTITIALS_CORE_BASE_SAFE_BROWSING_ERROR_UI_H_

#define CanShowEnhancedProtectionMessage               \
  CanShowEnhancedProtectionMessage() { return false; } \
  bool CanShowEnhancedProtectionMessage_ChromiumImpl

#include "src/components/security_interstitials/core/base_safe_browsing_error_ui.h"
#undef CanShowEnhancedProtectionMessage

#endif  // BEACON_CHROMIUM_SRC_COMPONENTS_SECURITY_INTERSTITIALS_CORE_BASE_SAFE_BROWSING_ERROR_UI_H_
