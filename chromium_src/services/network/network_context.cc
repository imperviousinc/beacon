#include "beacon/services/trust/client/dnssec_cert_verifier_factory.h"

// Wraps the Web PKI verifier with our DANE implementation
// proper mojo remote should be plumbed here
#define BEACON_WRAP_WEB_PKI_DNSSEC_VERIFIER \
    cert_verifier = DNSSECCertVerifierFactory::Create(std::move(cert_verifier));

#include "src/services/network/network_context.cc"
#undef BEACON_WRAP_WEB_PKI_DNSSEC_VERIFIER
