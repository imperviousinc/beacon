#ifndef BEACON_SERVICES_TRUST_CLIENT_DNSSEC_CERT_VERIFIER_FACTORY_H_
#define BEACON_SERVICES_TRUST_CLIENT_DNSSEC_CERT_VERIFIER_FACTORY_H_

#include "net/cert/cert_verifier.h"
#include "beacon/services/trust/client/dnssec_cert_verifier.h"

class DNSSECCertVerifierFactory {
    public:
    static std::unique_ptr<DNSSECCertVerifier> Create(std::unique_ptr<net::CertVerifier> upstream);
};


#endif // BEACON_SERVICES_TRUST_CLIENT_DNSSEC_CERT_VERIFIER_FACTORY_H_
