#define BEACON_ADD_PAGE_INFO_DETAILS                                               \
if (identity_info.certificate->is_dnssec_cert) {                                   \
  auto description = CreateSecurityDescription(                                    \
          SecuritySummaryColor::GREEN, IDS_BEACON_PAGE_INFO_DNSSEC_SECURE_SUMMARY, \
          IDS_BEACON_PAGE_INFO_DNSSEC_SECURE_DETAILS,                              \
          SecurityDescriptionType::CONNECTION);                                    \
  return description;                                                              \
}

#include "src/components/page_info/page_info_ui.cc"
#undef BEACON_ADD_PAGE_INFO_DETAILS
