#include "src/net/base/net_error_list.h"

// Beacon net error codes start at 3500

// 3500 - 3599 DNSSEC errors
// 3600 - 3699 Light client/service errors

// DNSSEC validation failed ("bogus" state).
NET_ERROR(DNSSEC_BOGUS, -3500)

// An RRSIG used to validate some necessary RRSet 
// in the chain was expired.
NET_ERROR(DNSSEC_SIGNATURE_EXPIRED, -3501)

// A missing RRSIG is needed to validate some RRSet 
// in the chain.
NET_ERROR(DNSSEC_SIGNATURE_MISSING, -3502)

// Broken chain no valid path to a DNSKEY.
NET_ERROR(DNSSEC_DNSKEY_MISSING, -3503)

// Missing denial of existence proof (NSEC or NSEC3)
NET_ERROR(DNSSEC_NSEC_MISSING, -3504)

// A validated TLSA RRSet exists but no certificate
// matches any of the TLSA records.
NET_ERROR(DNSSEC_PINNED_KEY_NOT_IN_CERT_CHAIN, -3505)

// A network error related to fetching
// the dnssec chain.
NET_ERROR(DNSSEC_FETCH_FAILED, -3506)

// Fetching DNSSEC chain timed out
NET_ERROR(DNSSEC_FETCH_TIMED_OUT, -3507)

// Still syncing.
NET_ERROR(HNS_IS_SYNCING, -3600)

// Looking for peers.
NET_ERROR(HNS_NO_PEERS, -3601)

// Timed out requesting proofs from a peer.
NET_ERROR(HNS_PEER_TIMED_OUT, -3602)

// HNS light client could return a generic
// error for serveral reasons.
NET_ERROR(HNS_REQUEST_FAILED, -3603)

// Some HIP-5 handler needed to verify 
// this request timed out.
NET_ERROR(HNS_HIP5_HANDLER_TIMED_OUT, -3604)

// Some HIP-5 handler needed to verify 
// this request failed.
NET_ERROR(HNS_HIP5_HANDLER_FAILED, -3605)

// A request to the trust service failed (generic error)
NET_ERROR(TRUST_SERVICE_REQUEST_FAILED, -3606)

// A request to the trust service timed out
NET_ERROR(TRUST_SERVICE_REQUEST_TIMED_OUT, -3607)

// Sending a bad request to the trust service
NET_ERROR(TRUST_SERVICE_REQUEST_INVALID, -3608)

// Receiving a bad response from the trust service
NET_ERROR(TRUST_SERVICE_RESPONSE_INVALID, -3609)
