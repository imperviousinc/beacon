#ifndef BEACON_COMPONENTS_CORE_BINDINGS_BINDINGS_H_
#define BEACON_COMPONENTS_CORE_BINDINGS_BINDINGS_H_

#include "beacon/components/core/bindings/core_library.h"
#include "base/memory/weak_ptr.h"
#include "base/callback.h"

namespace beacon {
namespace core {

class Service {
  public:
    Service(std::unique_ptr<CoreLibrary> library);
    Service(const Service&) = delete;
    Service& operator=(const Service&) = delete;
    ~Service();

    // Launches Beacon core service in a dedicated thread
    // it will return false if already running
    bool Launch();

    // Sends a shutdown signal to the service
    void Shutdown();
  private:
    static int32_t BlockingLaunch(CoreLibrary::LaunchFunc func);
    void OnShutdown(int32_t code);

    bool started_ = false;
    bool stopping_ = false;
    // The task runner of the creation thread.
    scoped_refptr<base::SequencedTaskRunner> task_runner_;

    std::unique_ptr<CoreLibrary> core_;
    base::WeakPtrFactory<Service> weak_ptr_factory_{this};
};


} // namespace core
} // namespace beacon

#endif // BEACON_COMPONENTS_CORE_BINDINGS_BINDINGS_H_
