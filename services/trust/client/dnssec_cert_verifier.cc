
#include "beacon/services/trust/client/dnssec_cert_verifier.h"
#include "beacon/services/trust/client/grpc/grpc_client.h"
#include "net/base/registry_controlled_domains/registry_controlled_domain.h"
#include <memory>
#include "base/check.h"
#include "base/logging.h"
#include "net/base/url_util.h"
#include "net/cert/x509_util.h"
#include "net/cert/x509_certificate.h"
#include "net/cert/cert_verifier.h"
#include "net/cert/cert_verify_result.h"
#include "url/url_canon.h"
#include "beacon/services/trust/client/client_errors.h"

DNSSECCertVerifier::Request::Request() {}
DNSSECCertVerifier::Request::~Request() {}
DNSSECCertVerifier::~DNSSECCertVerifier() {}

void DNSSECCertVerifier::Request::OnRemoteResponse(
        const net::CertVerifier::RequestParams& params,
        net::CertVerifyResult* verify_result,
        int error_from_upstream,
        net::CompletionOnceCallback callback,
        const grpc::Status& grpcStatus,
        const dnssec_cert_verifier::CertVerifyResponse& response) {

  auto cert = params.certificate();

  if (!grpcStatus.ok()) {
    verify_result->Reset();
    verify_result->verified_cert = cert;
    verify_result->cert_status |= net::CERT_STATUS_INVALID;
    std::move(callback).Run(net::ERR_TRUST_SERVICE_REQUEST_FAILED);
    return;
  }

  auto state = response.state();
  if (state == dnssec_cert_verifier::SecurityState::INSECURE) {
    // This zone is insecure or doesn't have a TLSA record
    // pass upstream results to the callback.
    std::move(callback).Run(error_from_upstream);
    return;
  }

  verify_result->Reset();
  verify_result->verified_cert = cert;

  if (state == dnssec_cert_verifier::SecurityState::BOGUS) {
    verify_result->cert_status |= net::CERT_STATUS_INVALID;

    int net_error = beacon::core::MapClientErrorToNetError(response.code());
    CHECK(net_error != net::OK);

    std::move(callback).Run(net_error);
    return;
  }

  // Response must be secure
  CHECK(state == dnssec_cert_verifier::SECURE);

  scoped_refptr<net::X509Certificate> verified_cert = net::X509Certificate::CreateFromBuffer(
      bssl::UpRef(cert->cert_buffer()), {});

  // certificate should've been parsed successfully by upstream verifier
  CHECK(verified_cert);
  verified_cert->is_dnssec_cert = true;
  verified_cert->is_hns_hostname = true;

  verify_result->verified_cert = verified_cert;
  std::move(callback).Run(net::OK);
}

DNSSECCertVerifier::DNSSECCertVerifier(std::unique_ptr<net::CertVerifier> upstream, 
    std::shared_ptr<grpc::Channel> channel): upstream_(std::move(upstream)), channel_(channel) {
  
}

namespace beacon {
bool IsHostnameNonDNS(const std::string& hostname) {
  // CanonicalizeHost requires surrounding brackets to parse an IPv6 address.
  const std::string host_or_ip = hostname.find(':') != std::string::npos ?
      "[" + hostname + "]" : hostname;
  url::CanonHostInfo host_info;
  std::string canonical_name = net::CanonicalizeHost(host_or_ip, &host_info);

  // If canonicalization fails, then the input is truly malformed.
  if (canonical_name.empty())
    return true;

  // IP addresses aren't valid DNSSEC hostnames.  
  if (host_info.IsIPAddress())
    return true;   

  // Ignore ICANN TLDs
  return net::registry_controlled_domains::HostHasRegistryControlledDomain(
       canonical_name, net::registry_controlled_domains::EXCLUDE_UNKNOWN_REGISTRIES,
       net::registry_controlled_domains::EXCLUDE_PRIVATE_REGISTRIES);
}

} // namespace beacon

// CertVerifier implementation
int DNSSECCertVerifier::Verify(const net::CertVerifier::RequestParams& params,
             net::CertVerifyResult* verify_result,
             net::CompletionOnceCallback callback,
             std::unique_ptr<net::CertVerifier::Request>* out_req,
             const net::NetLogWithSource& net_log) {
    out_req->reset();

    // Pass Non-DNS hostnames callback to the upstream verifier 
    if (beacon::IsHostnameNonDNS(params.hostname())) {
        return upstream_->Verify(params, verify_result, std::move(callback), out_req, net_log);
    }

    // Pass our callback instead. We still call the upstream verifier
    // to potentially catch any serious errors before performing
    // our own verification.
    net::CompletionOnceCallback callback2 = base::BindOnce(
        &DNSSECCertVerifier::OnRequestFinished, base::Unretained(this),
        params, std::move(callback), verify_result, out_req);
        
    return upstream_->Verify(params, verify_result, std::move(callback2), out_req, net_log);
}


void DNSSECCertVerifier::SetConfig(const Config& config) {
    upstream_->SetConfig(config);
}

void DNSSECCertVerifier::OnRequestFinished(const net::CertVerifier::RequestParams& params,
                         net::CompletionOnceCallback callback,
                         net::CertVerifyResult* verify_result,
                         std::unique_ptr<net::CertVerifier::Request>* out_req,
                         int error) {
    // If the cert error was fatal we call
    // the original callback skipping 
    // DNSSEC verification.         
    if (error != net::OK && 
        error != net::ERR_CERT_AUTHORITY_INVALID && 
        error != net::ERR_CERT_DATE_INVALID && 
        error != net::ERR_CERT_COMMON_NAME_INVALID) {
        std::move(callback).Run(error);
        return;    
    }

    // Cert verify service call configuration.
    beacon::core::StateConfig state_config;
    state_config.max_retries = 5;
    state_config.timeout_in_ms = 10000; 
    state_config.wait_for_ready = true;

    LOG(INFO)<< "DNSSECCertVerifier on request received";

    dnssec_cert_verifier::CertVerifyRequest grpcRequest;
    grpcRequest.set_host(params.hostname());

    // Unfortunately, RequestParams doesn't provide port info
    // so we hardcode port 443 in the TLSA lookup for now.
    // Note: CachingCertVerifier uses the hash of RequestParams
    // as key so verifying the identity of this cert on this port
    // should allow the same cert/name to work for other ports.
    grpcRequest.set_port("443");

    // Packaging the cert for the verify request.
    auto cert = params.certificate();
    dnssec_cert_verifier::Certificate* grpcCert = grpcRequest.mutable_cert();

    auto leaf = net::x509_util::CryptoBufferAsStringPiece(cert->cert_buffer());
    grpcCert->add_der_certs(leaf.data(), leaf.size());

    for (const auto&buffer : cert->intermediate_buffers()) {
       auto inter = net::x509_util::CryptoBufferAsStringPiece(buffer.get());
       grpcCert->add_der_certs(inter.data(), inter.size());
    }

    // Take a weak pointer to the request because deletion of the request
    // is what signals cancellation. If the request is cancelled, the
    // callback won't be called, thus avoiding UAF, because |verify_result|
    // is freed when the request is cancelled.
    *out_req = std::make_unique<Request>();
    base::WeakPtr<Request> weak_req = static_cast<Request*>(out_req->get())->GetWeakPtr();

    beacon::core::GrpcClient beaconServiceClient = beacon::core::GrpcClient(channel_);
    beaconServiceClient.CallServiceMethod(grpcRequest, base::BindOnce(&Request::OnRemoteResponse, 
        weak_req, params, verify_result, error, std::move(callback)), state_config);
}
