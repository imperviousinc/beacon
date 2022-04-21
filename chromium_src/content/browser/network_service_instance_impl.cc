#include "beacon/content/browser/trust_service_instance_impl.h"
#include "beacon/services/trust/trust_service.h"
#include "beacon/services/trust/public/mojom/trust_service.mojom.h"

// This is called in GetNetworkService so it gets launched
// when the network service is launched.
#define BEACON_GET_TRUST_SERVICE beacon::content::GetTrustService();
// Shuts down TrustService when network service is shutdown
#define BEACON_SHUTDOWN_TRUST_SERVICE beacon::content::ShutDownTrustService();

#include "src/content/browser/network_service_instance_impl.cc"

#undef BEACON_GET_TRUST_SERVICE
#undef BEACON_SHUTDOWN_TRUST_SERVICE

