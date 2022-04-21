// Hooking Beacon Trust Service
#include "beacon/services/trust/public/mojom/trust_service.mojom.h" 
#include "beacon/services/trust/trust_service.h"  

auto RunTrustService(mojo::PendingReceiver<trust::mojom::TrustService> receiver) {
  return std::make_unique<trust::TrustService>(std::move(receiver));  
}


#define BEACON_REGISTER_UTILITY_SERVICES services.Add(RunTrustService);
#include "src/content/utility/services.cc"
#undef BEACON_REGISTER_UTILITY_SERVICES
