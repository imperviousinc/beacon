#include "components/security_interstitials/content/cert_report_helper.h"


#define BEACON_CERT_REPORT_HELPER_SHOULD_SHOW_ENHANCED_PROTECTION_MESSAGE \
  return false;
#include "src/components/security_interstitials/content/cert_report_helper.cc"
#undef BEACON_CERT_REPORT_HELPER_SHOULD_SHOW_ENHANCED_PROTECTION_MESSAGE