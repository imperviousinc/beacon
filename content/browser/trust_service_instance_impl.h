#ifndef BEACON_CONTENT_BROWSER_TRUST_SERVICE_INSTANCE_IMPL_H_
#define BEACON_CONTENT_BROWSER_TRUST_SERVICE_INSTANCE_IMPL_H_

#include "base/callback.h"
#include "base/callback_list.h"
#include "content/common/content_export.h"
#include "beacon/services/trust/trust_service.h"
#include "beacon/services/trust/public/mojom/trust_service.mojom.h"

namespace beacon {
namespace content {

const scoped_refptr<base::SequencedTaskRunner>& GetTrustTaskRunner();

// Returns a pointer to the TrustService, creating / re-creating it as needed.
// This method can only be called on the UI thread.
trust::mojom::TrustService* GetTrustService();

// Registers |handler| to run (on UI thread) after mojo::Remote<TrustService>
// encounters an error. 
//
// Can only be called on the UI thread. 
base::CallbackListSubscription
RegisterTrustServiceCrashHandler(base::RepeatingClosure handler);

// Shuts down the in-process trust service or disconnects from the out-of-
// process one, allowing it to shut down.
void ShutDownTrustService();

} // namespace content
} // namespace beacon

#endif  // BEACON_CONTENT_BROWSER_TRUST_SERVICE_INSTANCE_IMPL_H_
