#include "beacon/services/trust/client/client_errors.h"
#include "beacon/services/trust/client/proto/dnssec_cert_verifier.grpc.pb.h"
#include "net/base/net_errors.h"

using dnssec_cert_verifier::ErrorCode;

namespace beacon {
namespace core {

int MapClientErrorToNetError(int error) {
  switch(error) {
  // DNSSEC errors
  case ErrorCode::ERR_DNSSEC_BOGUS:
    return net::ERR_DNSSEC_BOGUS;
  case ErrorCode::ERR_DNSSEC_SIGNATURE_EXPIRED:
    return net::ERR_DNSSEC_SIGNATURE_EXPIRED;
  case ErrorCode::ERR_DNSSEC_SIGNATURE_MISSING:
    return net::ERR_DNSSEC_SIGNATURE_MISSING;
  case ErrorCode::ERR_DNSSEC_DNSKEY_MISSING:
    return net::ERR_DNSSEC_DNSKEY_MISSING;
  case ErrorCode::ERR_DNSSEC_NSEC_MISSING:
    return net::ERR_DNSSEC_NSEC_MISSING;
  case ErrorCode::ERR_DNSSEC_PINNED_KEY_NOT_IN_CERT_CHAIN:
    return net::ERR_DNSSEC_PINNED_KEY_NOT_IN_CERT_CHAIN;
  case ErrorCode::ERR_DNSSEC_FETCH_FAILED:
    return net::ERR_DNSSEC_FETCH_FAILED;
  case ErrorCode::ERR_DNSSEC_FETCH_TIMED_OUT:
    return net::ERR_DNSSEC_FETCH_TIMED_OUT;
  // HNS errors
  case ErrorCode::ERR_HNS_IS_SYNCING:
    return net::ERR_HNS_IS_SYNCING;
  case ErrorCode::ERR_HNS_NO_PEERS:
    return net::ERR_HNS_NO_PEERS;
  case ErrorCode::ERR_HNS_PEER_TIMED_OUT:
    return net::ERR_HNS_PEER_TIMED_OUT;
  case ErrorCode::ERR_HNS_REQUEST_FAILED:
    return net::ERR_HNS_REQUEST_FAILED;
  case ErrorCode::ERR_HNS_HIP5_HANDLER_TIMED_OUT:
    return net::ERR_HNS_HIP5_HANDLER_TIMED_OUT;
  case ErrorCode::ERR_HNS_HIP5_HANDLER_FAILED:
    return net::ERR_HNS_HIP5_HANDLER_FAILED;
  // Trust service communication errors
  case ErrorCode::ERR_TRUST_SERVICE_REQUEST_FAILED:
    return net::ERR_TRUST_SERVICE_REQUEST_FAILED;
  case ErrorCode::ERR_TRUST_SERVICE_REQUEST_TIMED_OUT:
    return net::ERR_TRUST_SERVICE_REQUEST_TIMED_OUT;
  case ErrorCode::ERR_TRUST_SERVICE_REQUEST_INVALID:
    return net::ERR_TRUST_SERVICE_REQUEST_INVALID;
  case ErrorCode::ERR_TRUST_SERVICE_RESPONSE_INVALID:
    return net::ERR_TRUST_SERVICE_RESPONSE_INVALID;
  // Some net errors from chromium
  case ErrorCode::ERR_DNS_SECURE_RESOLVER_HOSTNAME_RESOLUTION_FAILED:
    return net::ERR_DNS_SECURE_RESOLVER_HOSTNAME_RESOLUTION_FAILED;
  case ErrorCode::ERR_DNS_TIMED_OUT:
    return net::ERR_DNS_TIMED_OUT;
  case ErrorCode::ERR_DNS_SERVER_FAILED:
    return net::ERR_DNS_SERVER_FAILED;
  case ErrorCode::ERR_DNS_MALFORMED_RESPONSE:
    return net::ERR_DNS_MALFORMED_RESPONSE;
  case ErrorCode::ERR_DNS_REQUEST_CANCELLED:
    return net::ERR_DNS_REQUEST_CANCELLED;
  // Cert errors from chromium
  case ErrorCode::ERR_CERT_COMMON_NAME_INVALID:
    return net::ERR_CERT_COMMON_NAME_INVALID;
  case ErrorCode::ERR_CERT_DATE_INVALID:
    return net::ERR_CERT_DATE_INVALID;
  case ErrorCode::ERR_CERT_AUTHORITY_INVALID:
    return net::ERR_CERT_AUTHORITY_INVALID;
  case ErrorCode::ERR_CERT_REVOKED:
    return net::ERR_CERT_REVOKED;
  case ErrorCode::ERR_CERT_INVALID:
    return net::ERR_CERT_INVALID;
  case ErrorCode::ERR_CERT_WEAK_SIGNATURE_ALGORITHM:
    return net::ERR_CERT_WEAK_SIGNATURE_ALGORITHM;
  case ErrorCode::ERR_CERT_NON_UNIQUE_NAME:
    return net::ERR_CERT_NON_UNIQUE_NAME;
  case ErrorCode::ERR_CERT_WEAK_KEY:
    return net::ERR_CERT_WEAK_KEY;
  case ErrorCode::ERR_CERT_NAME_CONSTRAINT_VIOLATION:
    return net::ERR_CERT_NAME_CONSTRAINT_VIOLATION;
  case ErrorCode::ERR_CERT_VALIDITY_TOO_LONG:
    return net::ERR_CERT_VALIDITY_TOO_LONG;
  case ErrorCode::ERR_CERT_KNOWN_INTERCEPTION_BLOCKED:
    return net::ERR_CERT_KNOWN_INTERCEPTION_BLOCKED;
  // Some other generic errors
  case ErrorCode::ERR_FAILED:
    return net::ERR_FAILED;
  case ErrorCode::ERR_ABORTED:
    return net::ERR_ABORTED;
  default:
    return net::ERR_UNEXPECTED;
  }
}

} // namespace trust
} // namespace beacon
