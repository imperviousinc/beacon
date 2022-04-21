#ifndef BEACON_SERVICES_TRUST_CLIENT_DNSSEC_CERT_VERIFIER_H_
#define BEACON_SERVICES_TRUST_CLIENT_DNSSEC_CERT_VERIFIER_H_

#include "net/cert/cert_verifier.h"
#include "net/cert/cert_verify_result.h"
#include "beacon/services/trust/client/proto/dnssec_cert_verifier.grpc.pb.h"
#include <grpcpp/channel.h>

class DNSSECCertVerifier : public net::CertVerifier {
 public:
  class Request : public net::CertVerifier::Request {
   public:
    Request();
    ~Request() override;
    void OnRemoteResponse(
        const net::CertVerifier::RequestParams& params,
        net::CertVerifyResult* verify_result,
        int error_from_upstream,
        net::CompletionOnceCallback callback,
        const grpc::Status& grpcStatus,
        const dnssec_cert_verifier::CertVerifyResponse& response);
    base::WeakPtr<Request> GetWeakPtr() { return weak_factory_.GetWeakPtr(); }
   private:
    base::WeakPtrFactory<Request> weak_factory_{this};
  };

  DNSSECCertVerifier(std::unique_ptr<net::CertVerifier> upstream, std::shared_ptr<grpc::Channel> channel);
  ~DNSSECCertVerifier() override;

  // CertVerifier implementation
  int Verify(const net::CertVerifier::RequestParams& params,
             net::CertVerifyResult* verify_result,
             net::CompletionOnceCallback callback,
             std::unique_ptr<net::CertVerifier::Request>* out_req,
             const net::NetLogWithSource& net_log) override;


  void SetConfig(const Config& config) override;

  void OnRequestFinished(const net::CertVerifier::RequestParams& params,
                         net::CompletionOnceCallback callback,
                         net::CertVerifyResult* verify_result,
                         std::unique_ptr<net::CertVerifier::Request>* out_req,
                         int error);

 private:
  std::unique_ptr<net::CertVerifier> upstream_;
  std::shared_ptr<grpc::Channel> channel_;
};

#endif // BEACON_SERVICES_TRUST_CLIENT_DNSSEC_CERT_VERIFIER_H_
