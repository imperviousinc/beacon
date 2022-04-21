#ifndef BEACON_SERVICES_TRUST_TRUST_SERVICE_H_
#define BEACON_SERVICES_TRUST_TRUST_SERVICE_H_

#include "beacon/services/trust/public/mojom/trust_service.mojom.h"
#include "mojo/public/cpp/bindings/receiver.h"
#include "beacon/components/core/bindings/bindings.h"

namespace trust {

class TrustService : public mojom::TrustService {
 public:
  explicit TrustService(mojo::PendingReceiver<mojom::TrustService> receiver);
  TrustService(const TrustService&) = delete;
  TrustService& operator=(const TrustService&) = delete;
  ~TrustService() override;

 private:
  // mojom::TrustService:
  void Launch() override;
  // TODO: once grpc is completely removed
  // define mojo bindings for cert verifier here
  mojo::Receiver<mojom::TrustService> receiver_;
  std::unique_ptr<beacon::core::Service> service_;
};

}  // namespace trust

#endif // BEACON_SERVICES_TRUST_TRUST_SERVICE_H_
