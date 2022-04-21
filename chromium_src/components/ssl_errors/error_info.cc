#include "components/ssl_errors/error_info.h"
#include "net/base/net_errors.h"

ssl_errors::ErrorInfo::ErrorType NetErrorToErrorType_BeaconImpl(int net_error) {
  switch (net_error) {
    case net::ERR_DNSSEC_BOGUS:
    case net::ERR_DNSSEC_NSEC_MISSING:
    case net::ERR_DNSSEC_DNSKEY_MISSING:
    case net::ERR_DNSSEC_SIGNATURE_EXPIRED:
    case net::ERR_DNSSEC_SIGNATURE_MISSING:
      return ssl_errors::ErrorInfo::CERT_INVALID;
    case net::ERR_DNSSEC_PINNED_KEY_NOT_IN_CERT_CHAIN:
      return ssl_errors::ErrorInfo::CERT_PINNED_KEY_MISSING;
    default:
      return ssl_errors::ErrorInfo::UNKNOWN;
  }
}

#define BEACON_TRY_CUSTOM_ERROR_TYPES { \
  auto error_type = NetErrorToErrorType_BeaconImpl(net_error); \
  if (error_type != ssl_errors::ErrorInfo::UNKNOWN) return error_type; \
}

#include "src/components/ssl_errors/error_info.cc"
#undef BEACON_TRY_CUSTOM_ERROR_TYPES

