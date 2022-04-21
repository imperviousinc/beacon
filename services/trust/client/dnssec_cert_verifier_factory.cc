#include "beacon/services/trust/client/dnssec_cert_verifier_factory.h"
#include "beacon/services/trust/client/dnssec_cert_verifier.h"
#include "net/cert/x509_util.h"
#include "net/cert/x509_certificate.h"
#include <grpcpp/create_channel.h>
#include "net/cert/cert_verifier.h"

// static
std::unique_ptr<DNSSECCertVerifier> DNSSECCertVerifierFactory::Create(std::unique_ptr<net::CertVerifier> upstream) {
    // TODO: replace with mojo once plumbing is done.
    auto remote = std::shared_ptr<grpc::Channel>(grpc::CreateChannel("127.0.0.1:44961", 
        grpc::InsecureChannelCredentials()));  

    return std::make_unique<DNSSECCertVerifier>(std::move(upstream), remote);
}