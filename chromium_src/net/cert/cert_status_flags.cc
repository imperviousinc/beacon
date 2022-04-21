#define MapNetErrorToCertStatus MapNetErrorToCertStatus_ChromiumImpl
#include "src/net/cert/cert_status_flags.cc"
#undef MapNetErrorToCertStatus

namespace net {

CertStatus MapNetErrorToCertStatus(int error) {
  switch(error) {
    case ERR_DNSSEC_BOGUS:
    case ERR_DNSSEC_NSEC_MISSING:
    case ERR_DNSSEC_DNSKEY_MISSING:
    case ERR_DNSSEC_SIGNATURE_EXPIRED:
    case ERR_DNSSEC_SIGNATURE_MISSING:
      return CERT_STATUS_INVALID;
    case ERR_DNSSEC_PINNED_KEY_NOT_IN_CERT_CHAIN:
      return CERT_STATUS_PINNED_KEY_MISSING;
    default:
      return MapNetErrorToCertStatus_ChromiumImpl(error);
  }
}

}
