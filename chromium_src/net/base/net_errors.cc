#include "net/base/net_errors.h"

#define IsCertificateError IsCertificateError_ChromiumImpl
#include "src/net/base/net_errors.cc"
#undef IsCertificateError

namespace net {

bool IsCertificateError(int error) {
  if (net::IsCertificateError_ChromiumImpl(error))
    return true;

  switch(error) {
    case ERR_DNSSEC_BOGUS:
    case ERR_DNSSEC_NSEC_MISSING:
    case ERR_DNSSEC_DNSKEY_MISSING:
    case ERR_DNSSEC_SIGNATURE_EXPIRED:
    case ERR_DNSSEC_SIGNATURE_MISSING:
    case ERR_DNSSEC_PINNED_KEY_NOT_IN_CERT_CHAIN:
      return true;
    default:
      return false;
  }
}

}
