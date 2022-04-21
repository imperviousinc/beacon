#include "net/base/net_errors.h"

int IsCertErrorFatal_BeaconImpl(int cert_error) {
  switch(cert_error) {
    case net::ERR_DNSSEC_BOGUS:
    case net::ERR_DNSSEC_NSEC_MISSING:
    case net::ERR_DNSSEC_DNSKEY_MISSING:
    case net::ERR_DNSSEC_SIGNATURE_EXPIRED:
    case net::ERR_DNSSEC_SIGNATURE_MISSING:
    case net::ERR_DNSSEC_PINNED_KEY_NOT_IN_CERT_CHAIN:
      return true;
    default:
      return -1;    
  }
}

#define BEACON_IS_CERT_ERROR_FATAL {\
  int res = IsCertErrorFatal_BeaconImpl(cert_error); \
  if (res != -1) return res; \
}

#include "src/components/security_interstitials/core/ssl_error_options_mask.cc"
#undef BEACON_IS_CERT_ERROR_FATAL
