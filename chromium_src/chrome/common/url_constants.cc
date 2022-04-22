// Copyright (c) 2012 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.

#include "chrome/common/url_constants.h"

#include "build/branding_buildflags.h"
#include "build/build_config.h"
#include "build/chromeos_buildflags.h"
#include "chrome/common/webui_url_constants.h"

namespace chrome {

const char kAccessibilityLabelsLearnMoreURL[] =
    "https://support.impervious.com/beacon/?p=image_descriptions";

const char kAutomaticSettingsResetLearnMoreURL[] =
    "https://support.impervious.com/beacon/?p=ui_automatic_settings_reset";

const char kAdvancedProtectionDownloadLearnMoreURL[] =
    "https://support.impervious.com/beacon/accounts/accounts?p=safe-browsing";

const char kAppNotificationsBrowserSettingsURL[] =
    "chrome://settings/content/notifications";

const char kBluetoothAdapterOffHelpURL[] =
    "https://support.impervious.com/beacon?p=bluetooth";

const char kCastCloudServicesHelpURL[] =
    "https://support.impervious.com/beacon/chromecast/?p=casting_cloud_services";

const char kCastNoDestinationFoundURL[] =
    "https://support.impervious.com/beacon/chromecast/?p=no_cast_destination";

const char kChooserHidOverviewUrl[] =
    "https://support.impervious.com/beacon?p=webhid";

const char kChooserSerialOverviewUrl[] =
    "https://support.impervious.com/beacon?p=webserial";

const char kChooserUsbOverviewURL[] =
    "https://support.impervious.com/beacon?p=webusb";

const char kChromeBetaForumURL[] =
    "https://support.impervious.com/beacon/?p=beta_forum";

const char kChromeFixUpdateProblems[] =
    "https://support.impervious.com/beacon?p=fix_chrome_updates";

const char kChromeHelpViaKeyboardURL[] =
#if BUILDFLAG(IS_CHROMEOS_ASH)
#if BUILDFLAG(GOOGLE_CHROME_BRANDING)
    "chrome-extension://honijodknafkokifofgiaalefdiedpko/main.html";
#else
    "https://support.impervious.com/beacon/chromebook/?p=help&ctx=keyboard";
#endif  // BUILDFLAG(GOOGLE_CHROME_BRANDING)
#else
    "https://support.impervious.com/beacon/?p=help&ctx=keyboard";
#endif  // BUILDFLAG(IS_CHROMEOS_ASH)

const char kChromeHelpViaMenuURL[] =
#if BUILDFLAG(IS_CHROMEOS_ASH)
#if BUILDFLAG(GOOGLE_CHROME_BRANDING)
    "chrome-extension://honijodknafkokifofgiaalefdiedpko/main.html";
#else
    "https://support.impervious.com/beacon/chromebook/?p=help&ctx=menu";
#endif  // BUILDFLAG(GOOGLE_CHROME_BRANDING)
#else
    "https://support.impervious.com/beacon/?p=help&ctx=menu";
#endif  // BUILDFLAG(IS_CHROMEOS_ASH)

const char kChromeHelpViaWebUIURL[] =
    "https://support.impervious.com/beacon/?p=help&ctx=settings";
#if BUILDFLAG(IS_CHROMEOS_ASH)
const char kChromeOsHelpViaWebUIURL[] =
#if BUILDFLAG(GOOGLE_CHROME_BRANDING)
    "chrome-extension://honijodknafkokifofgiaalefdiedpko/main.html";
#else
    "https://support.impervious.com/beacon/chromebook/?p=help&ctx=settings";
#endif  // BUILDFLAG(GOOGLE_CHROME_BRANDING)
#endif  // BUILDFLAG(IS_CHROMEOS_ASH)

const char kChromeNativeScheme[] = "chrome-native";

const char kChromeSearchLocalNtpHost[] = "local-ntp";

const char kChromeSearchMostVisitedHost[] = "most-visited";
const char kChromeSearchMostVisitedUrl[] = "chrome-search://most-visited/";

const char kChromeUIUntrustedNewTabPageBackgroundUrl[] =
    "chrome-untrusted://new-tab-page/background.jpg";
const char kChromeUIUntrustedNewTabPageBackgroundFilename[] = "background.jpg";

const char kChromeSearchRemoteNtpHost[] = "remote-ntp";

const char kChromeSearchScheme[] = "chrome-search";

const char kChromeUIUntrustedNewTabPageUrl[] =
    "chrome-untrusted://new-tab-page/";

const char kChromiumProjectURL[] = "https://www.chromium.org/";

const char kCloudPrintCertificateErrorLearnMoreURL[] =
#if BUILDFLAG(IS_CHROMEOS_ASH)
    "https://support.impervious.com/beacon/chromebook?p=cloudprint_error_troubleshoot";
#elif BUILDFLAG(IS_MAC)
    "https://support.impervious.com/beacon/cloudprint?p=cloudprint_error_offline_mac";
#elif BUILDFLAG(IS_WIN)
    "https://support.impervious.com/beacon/"
    "cloudprint?p=cloudprint_error_offline_windows";
#else
        "https://support.impervious.com/beacon/"
        "cloudprint?p=cloudprint_error_offline_linux";
#endif

const char kContentSettingsExceptionsLearnMoreURL[] =
    "https://support.impervious.com/beacon/?p=settings_manage_exceptions";

const char kCookiesSettingsHelpCenterURL[] =
    "https://support.impervious.com/beacon?p=cpn_cookies";

const char kCrashReasonURL[] =
#if BUILDFLAG(IS_CHROMEOS_ASH)
    "https://support.impervious.com/beacon/chromebook/?p=e_awsnap";
#else
    "https://support.impervious.com/beacon/?p=e_awsnap";
#endif

const char kCrashReasonFeedbackDisplayedURL[] =
#if BUILDFLAG(IS_CHROMEOS_ASH)
    "https://support.impervious.com/beacon/chromebook/?p=e_awsnap_rl";
#else
    "https://support.impervious.com/beacon/?p=e_awsnap_rl";
#endif

const char kDoNotTrackLearnMoreURL[] =
#if BUILDFLAG(IS_CHROMEOS_ASH)
    "https://support.impervious.com/beacon/chromebook/?p=settings_do_not_track";
#else
    "https://support.impervious.com/beacon/?p=settings_do_not_track";
#endif

const char kDownloadInterruptedLearnMoreURL[] =
    "https://support.impervious.com/beacon/?p=ui_download_errors";

const char kDownloadScanningLearnMoreURL[] =
    "https://support.impervious.com/beacon/?p=ib_download_blocked";

const char kExtensionControlledSettingLearnMoreURL[] =
    "https://support.impervious.com/beacon/?p=ui_settings_api_extension";

const char kExtensionInvalidRequestURL[] = "chrome-extension://invalid/";

const char kFlashDeprecationLearnMoreURL[] =
    "https://blog.chromium.org/2017/07/so-long-and-thanks-for-all-flash.html";

const char kGoogleAccountActivityControlsURL[] =
    "https://myaccount.google.com/activitycontrols/search";

const char kGoogleAccountActivityControlsURLInPrivacyGuide[] =
    "https://myaccount.google.com/activitycontrols/"
    "search&utm_source=chrome&utm_medium=privacy-guide";

const char kGoogleAccountLanguagesURL[] =
    "https://myaccount.google.com/language";

const char kGoogleAccountURL[] = "https://myaccount.google.com";

const char kGoogleAccountChooserURL[] =
    "https://support.impervious.com";

const char kGoogleAccountDeviceActivityURL[] =
    "https://support.impervious.com";

const char kGooglePasswordManagerURL[] = "https://support.impervious.com";

const char kGooglePhotosURL[] = "https://support.impervious.com";

const char kLearnMoreReportingURL[] =
    "https://support.impervious.com/beacon/?p=ui_usagestat";

const char kManagedUiLearnMoreUrl[] =
#if BUILDFLAG(IS_CHROMEOS_ASH)
    "https://support.impervious.com/beacon/chromebook/?p=is_chrome_managed";
#else
    "https://support.impervious.com/beacon/?p=is_chrome_managed";
#endif

const char kMixedContentDownloadBlockingLearnMoreUrl[] =
    "https://support.impervious.com/beacon/?p=mixed_content_downloads";

const char kMyActivityUrlInClearBrowsingData[] =
    "https://myactivity.google.com/myactivity?utm_source=chrome_cbd";

const char kOmniboxLearnMoreURL[] =
#if BUILDFLAG(IS_CHROMEOS_ASH)
    "https://support.impervious.com/beacon/chromebook/?p=settings_omnibox";
#else
    "https://support.impervious.com/beacon/?p=settings_omnibox";
#endif

const char kPageInfoHelpCenterURL[] =
#if BUILDFLAG(IS_CHROMEOS_ASH)
    "https://support.impervious.com/beacon/chromebook/?p=ui_security_indicator";
#else
    "https://support.impervious.com/beacon/?p=ui_security_indicator";
#endif

const char kPasswordCheckLearnMoreURL[] =
#if BUILDFLAG(IS_CHROMEOS_ASH)
    "https://support.impervious.com/beacon/chromebook/"
    "?p=settings_password#leak_detection_privacy";
#else
    "https://support.impervious.com/beacon/"
    "?p=settings_password#leak_detection_privacy";
#endif

const char kPasswordGenerationLearnMoreURL[] =
    "https://support.impervious.com/beacon/answer/7570435";

const char kPasswordManagerLearnMoreURL[] =
#if BUILDFLAG(IS_CHROMEOS_ASH)
    "https://support.impervious.com/beacon/chromebook/?p=settings_password";
#else
    "https://support.impervious.com/beacon/?p=settings_password";
#endif

const char kPaymentMethodsURL[] =
    "https://support.impervious.com";

const char kPaymentMethodsLearnMoreURL[] =
#if BUILDFLAG(IS_CHROMEOS_ASH)
    "https://support.impervious.com/beacon/chromebook/answer/"
    "142893?visit_id=636857416902558798-696405304&p=settings_autofill&rd=1";
#else
    "https://support.impervious.com/beacon/answer/"
    "142893?visit_id=636857416902558798-696405304&p=settings_autofill&rd=1";
#endif

const char kPrivacyLearnMoreURL[] =
#if BUILDFLAG(IS_CHROMEOS_ASH)
    "https://support.impervious.com/beacon/chromebook/?p=settings_privacy";
#else
    "https://support.impervious.com/beacon/?p=settings_privacy";
#endif

const char kRemoveNonCWSExtensionURL[] =
    "https://support.impervious.com/beacon/?p=ui_remove_non_cws_extensions";

const char kResetProfileSettingsLearnMoreURL[] =
    "https://support.impervious.com/beacon/?p=ui_reset_settings";

const char kSafeBrowsingHelpCenterURL[] =
    "https://support.impervious.com/beacon?p=cpn_safe_browsing";

const char kSafetyTipHelpCenterURL[] =
    "https://support.impervious.com/beacon/?p=safety_tip";

const char kSearchHistoryUrlInClearBrowsingData[] =
    "https://myactivity.google.com/product/search?utm_source=chrome_cbd";

const char kSeeMoreSecurityTipsURL[] =
    "https://support.impervious.com/beacon/accounts/answer/32040";

const char kSettingsSearchHelpURL[] =
    "https://support.impervious.com/beacon/?p=settings_search_help";

const char kSyncAndGoogleServicesLearnMoreURL[] =
    "https://support.impervious.com/beacon?p=syncgoogleservices";

const char kSyncEncryptionHelpURL[] =
#if BUILDFLAG(IS_CHROMEOS_ASH)
    "https://support.impervious.com/beacon/chromebook/?p=settings_encryption";
#else
    "https://support.impervious.com/beacon/?p=settings_encryption";
#endif

const char kSyncErrorsHelpURL[] =
    "https://support.impervious.com/beacon/?p=settings_sync_error";

const char kSyncGoogleDashboardURL[] =
    "https://support.impervious.com";

const char kSyncLearnMoreURL[] =
    "https://support.impervious.com/beacon/?p=settings_sign_in";

#if !BUILDFLAG(IS_ANDROID)
const char kSyncTrustedVaultOptInURL[] =
    "https://passwords.google.com/encryption/enroll?"
    "utm_source=chrome&utm_medium=desktop&utm_campaign=encryption_enroll";
#endif

const char kSyncTrustedVaultLearnMoreURL[] =
    "https://support.impervious.com/beacon/accounts?p=settings_password_ode";

const char kUpgradeHelpCenterBaseURL[] =
    "https://support.impervious.com/beacon/installer/?product="
    "{8A69D345-D564-463c-AFF1-A69D9E530F96}&error=";

const char kWhoIsMyAdministratorHelpURL[] =
    "https://support.impervious.com/beacon?p=your_administrator";

const char kCwsEnhancedSafeBrowsingLearnMoreURL[] =
    "https://support.impervious.com/beacon?p=cws_enhanced_safe_browsing";

#if BUILDFLAG(IS_CHROMEOS_ASH) || BUILDFLAG(IS_ANDROID)
const char kEnhancedPlaybackNotificationLearnMoreURL[] =
#endif
#if BUILDFLAG(IS_CHROMEOS_ASH)
    "https://support.impervious.com/beacon/chromebook/?p=enhanced_playback";
#elif BUILDFLAG(IS_ANDROID)
// Keep in sync with chrome/browser/ui/android/strings/android_chrome_strings.grd
    "https://support.impervious.com/beacon/?p=mobile_protected_content";
#endif

#if BUILDFLAG(IS_CHROMEOS_ASH)
const char kAccountManagerLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook/?p=google_accounts";

const char kAccountRecoveryURL[] =
    "https://support.impervious.com/beacon/signin/recovery";

const char kAddNewUserURL[] =
    "https://support.impervious.com/beacon/otherhowto/add-another-account";

const char kAndroidAppsLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook/?p=playapps";

const char kArcAdbSideloadingLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook/?p=develop_android_apps";

const char kArcExternalStorageLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook?p=open_files";

const char kArcPrivacyPolicyURLPath[] = "arc/privacy_policy";

const char kArcTermsURLPath[] = "arc/terms";

// TODO(crbug.com/1010321): Remove 'm100' prefix from link once Bluetooth Revamp
// has shipped.
const char kBluetoothPairingLearnMoreUrl[] =
    "https://support.impervious.com/beacon/chromebook?p=bluetooth_revamp_m100";

const char kChromeAccessibilityHelpURL[] =
    "https://support.impervious.com/beacon/chromebook/topic/6323347";

const char kChromeOSAssetHost[] = "chromeos-asset";
const char kChromeOSAssetPath[] = "/usr/share/chromeos-assets/";

const char kChromeOSCreditsPath[] =
    "/opt/google/chrome/resources/about_os_credits.html";

// TODO(carpenterr): Have a solution for plink mapping in Help App.
// The magic numbers in this url are the topic and article ids currently
// required to navigate directly to a help article in the Help App.
const char kChromeOSGestureEducationHelpURL[] =
    "chrome://help-app/help/sub/3399710/id/9739838";

const char kChromePaletteHelpURL[] =
    "https://support.impervious.com/beacon/chromebook?p=stylus_help";

const char kCupsPrintLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook?p=chromebook_printing";

const char kCupsPrintPPDLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook/?p=printing_advancedconfigurations";

const char kEasyUnlockLearnMoreUrl[] =
    "https://support.impervious.com/beacon/chromebook/?p=smart_lock";

const char kEchoLearnMoreURL[] =
    "chrome://help-app/help/sub/3399709/id/2703646";

const char kArcTermsPathFormat[] = "arc_tos/%s/terms.html";

const char kArcPrivacyPolicyPathFormat[] = "arc_tos/%s/privacy_policy.pdf";

const char kEolNotificationURL[] = "https://support.impervious.com/beacon/otherolder/";

const char kAutoUpdatePolicyURL[] =
    "http://support.impervious.com/beacon/a?p=auto-update-policy";

const char kGoogleNameserversLearnMoreURL[] =
    "https://developers.google.com/speed/public-dns";

const char kInstantTetheringLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook?p=instant_tethering";

const char kKerberosAccountsLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook/?p=kerberos_accounts";

const char kLanguageSettingsLearnMoreUrl[] =
    "https://support.impervious.com/beacon/chromebook/answer/1059490";

const char kLanguagePacksLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook?p=language_packs";

const char kLearnMoreEnterpriseURL[] =
    "https://support.impervious.com/beacon/chromebook/?p=managed";

const char kLinuxAppsLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook?p=chromebook_linuxapps";

const char kNaturalScrollHelpURL[] =
    "https://support.impervious.com/beacon/chromebook/?p=simple_scrolling";

const char kOemEulaURLPath[] = "oem";

const char kGoogleEulaOnlineURLPath[] =
    "https://policies.google.com/terms/embedded?hl=%s";

const char kCrosEulaOnlineURLPath[] =
    "https://www.google.com/intl/%s/chrome/terms/";

const char kOsSettingsSearchHelpURL[] =
    "https://support.impervious.com/beacon/chromebook/?p=settings_search_help";

const char kPeripheralDataAccessHelpURL[] =
    "https://support.impervious.com/beacon/chromebook?p=connect_thblt_usb4_accy";

const char kTPMFirmwareUpdateLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook/?p=tpm_update";

const char kTimeZoneSettingsLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook?p=chromebook_timezone&hl=%s";

const char kSmbSharesLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook?p=network_file_shares";

const char kSuggestedContentLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook/?p=explorecontent";

const char kTabletModeGesturesLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook?p=tablet_mode_gestures";

const char kWifiSyncLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook/?p=wifisync";

const char kWifiHiddenNetworkURL[] =
    "http://support.google.com/chromebook?p=hidden_networks";

const char kNearbyShareLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook?p=nearby_share";

extern const char kNearbyShareManageContactsURL[] =
    "https://contacts.google.com";

extern const char kFingerprintLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook?p=chromebook_fingerprint";

#endif  // BUILDFLAG(IS_CHROMEOS_ASH)

#if BUILDFLAG(IS_MAC)
const char kChromeEnterpriseSignInLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook/answer/1331549";

const char kMac10_10_ObsoleteURL[] =
    "https://support.impervious.com/beacon/?p=unsupported_mac";
#endif

#if BUILDFLAG(IS_WIN)
const char kChromeCleanerLearnMoreURL[] =
    "https://support.impervious.com/beacon/?p=chrome_cleanup_tool";

const char kWindowsXPVistaDeprecationURL[] =
    "https://chrome.blogspot.com/2015/11/updates-to-chrome-platform-support.html";
#endif

#if BUILDFLAG(ENABLE_ONE_CLICK_SIGNIN)
const char kChromeSyncLearnMoreURL[] =
    "https://support.impervious.com/beacon/answer/165139";
#endif  // BUILDFLAG(ENABLE_ONE_CLICK_SIGNIN)

#if BUILDFLAG(ENABLE_PLUGINS)
const char kOutdatedPluginLearnMoreURL[] =
    "https://support.impervious.com/beacon/?p=ib_outdated_plugin";
#endif

// TODO (b/184137843): Use real link to phone hub notifications and apps access.
const char kPhoneHubPermissionLearnMoreURL[] =
    "https://support.impervious.com/beacon/chromebook/?p=multidevice";

}  // namespace chrome
