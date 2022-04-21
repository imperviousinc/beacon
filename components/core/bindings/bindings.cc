#include "beacon/components/core/bindings/bindings.h"
#include "beacon/components/core/bindings/core_library.h"
#include "base/threading/thread.h"
#include "base/task/post_task.h"
#include "base/callback.h"
#include "base/no_destructor.h"
#include "base/threading/scoped_blocking_call.h"
#include "base/logging.h"
#include "base/task/thread_pool.h"
#include "base/callback.h"
#include "base/threading/sequenced_task_runner_handle.h"
#include "base/task/bind_post_task.h"

namespace beacon {
namespace core {

base::Thread& GetCoreLibraryDedicatedThread() {
  static base::NoDestructor<base::Thread> thread{"BeaconCoreDedicatedThread"};
  return *thread;
}

Service::Service(std::unique_ptr<CoreLibrary> core) :
  task_runner_(base::SequencedTaskRunnerHandle::Get()), core_(std::move(core)) {
  CHECK(core_);
}

Service::~Service() = default;

// static
DISABLE_CFI_ICALL
int32_t Service::BlockingLaunch(CoreLibrary::LaunchFunc run) {
  base::ScopedBlockingCall scoped_blocking_call(FROM_HERE,
                                                base::BlockingType::WILL_BLOCK);
  LOG(INFO) << "Starting Beacon core service";
  int32_t code = run();
  return code;
}

void Service::Shutdown() {
  CHECK(task_runner_->RunsTasksInCurrentSequence());
  if (!started_ || stopping_) {
    LOG(INFO) << "Beacon service is already shutdown or is being shutdown";
    return;
  }

  stopping_ = true;
  core_->Shutdown();
}

bool Service::Launch() {
  CHECK(task_runner_->RunsTasksInCurrentSequence());

  if (started_ || stopping_) {
    LOG(INFO) << "Beacon service is already running";
    return false;
  }

  base::OnceCallback<void(int32_t)> shutdown_callback = base::BindOnce(&Service::OnShutdown,
           weak_ptr_factory_.GetWeakPtr());

  // Launch callback will be called by the dedicated thread task runner
  // we use a dedicated thread because this call will never return until
  // exit or on error. It is safe to call BeaconHelper_Launch from any thread
  // so we pass its pointer.
  auto launch_callback = base::BindOnce(&Service::BlockingLaunch, core_->GetLaunchFunc());

  // The task runner for this sequence will then call shutdown_callback
  // if |this| is still alive.
  base::OnceClosure task = std::move(launch_callback)
           .Then(base::BindPostTask(task_runner_, std::move(shutdown_callback)));

  if (!GetCoreLibraryDedicatedThread().IsRunning()) {
    base::Thread::Options options(base::MessagePumpType::IO, 0);
    GetCoreLibraryDedicatedThread().StartWithOptions(std::move(options));
  }

  auto dedicated_thread_task_runner = GetCoreLibraryDedicatedThread().task_runner();
  dedicated_thread_task_runner->PostTask(FROM_HERE, std::move(task));
  started_ = true;
  return true;
}

void Service::OnShutdown(int32_t code) {
  CHECK(task_runner_->RunsTasksInCurrentSequence());
  started_ = false;
  stopping_ = false;

  LOG(INFO) << "Beacon core service shutdown with code: " << code;
}


} // namespace core
} // namespace beacon
