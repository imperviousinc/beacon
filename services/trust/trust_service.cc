#include "beacon/services/trust/trust_service.h"
#include "beacon/components/core/bindings/bindings.h"
#include "beacon/components/core/bindings/core_library.h"

namespace trust {

TrustService::TrustService(mojo::PendingReceiver<mojom::TrustService> receiver)
    : receiver_(this, std::move(receiver)) {
  auto core_lib = beacon::core::CoreLibrary::Create();
  service_ = std::make_unique<beacon::core::Service>(std::move(core_lib));
}

TrustService::~TrustService() = default;

void TrustService::Launch() {
  service_->Launch();
}

}  // namespace trust
