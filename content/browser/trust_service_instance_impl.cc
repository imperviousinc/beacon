#include "beacon/content/browser/trust_service_instance_impl.h"

#include <memory>
#include <string>
#include <utility>

#include "base/bind.h"
#include "base/environment.h"
#include "base/feature_list.h"
#include "base/files/file.h"
#include "base/files/file_enumerator.h"
#include "base/files/file_util.h"
#include "base/message_loop/message_pump_type.h"
#include "base/metrics/field_trial_params.h"
#include "base/metrics/histogram_functions.h"
#include "base/metrics/histogram_macros.h"
#include "base/no_destructor.h"
#include "base/strings/string_util.h"
#include "base/strings/utf_string_conversions.h"
#include "base/synchronization/waitable_event.h"
#include "base/task/sequenced_task_runner.h"
#include "base/task/thread_pool.h"
#include "base/threading/sequence_local_storage_slot.h"
#include "base/threading/thread.h"
#include "base/threading/thread_restrictions.h"
#include "base/trace_event/trace_event.h"
#include "build/build_config.h"
#include "build/chromeos_buildflags.h"
#include "content/browser/browser_main_loop.h"
#include "content/public/browser/browser_task_traits.h"
#include "content/public/browser/browser_thread.h"
#include "content/public/browser/content_browser_client.h"
#include "content/public/browser/network_service_instance.h"
#include "content/public/browser/service_process_host.h"
#include "content/public/common/content_client.h"
#include "content/public/common/content_features.h"
#include "content/public/common/network_service_util.h"
#include "mojo/public/cpp/bindings/pending_receiver.h"
#include "mojo/public/cpp/bindings/remote.h"
#include "beacon/services/trust/public/mojom/trust_service.mojom.h"

#if BUILDFLAG(IS_ANDROID)
#include "content/common/android/cpu_affinity_setter.h"
#endif  // BUILDFLAG(IS_ANDROID)


namespace beacon {
namespace content {

namespace {

bool g_force_create_trust_service_directly = false;
mojo::Remote<trust::mojom::TrustService>* g_trust_service_remote =
    nullptr;
base::Time g_last_trust_service_crash;

base::Thread& GetTrustServiceDedicatedThread() {
  static base::NoDestructor<base::Thread> thread{"BeaconTrustService"};
  return *thread;
}

// The instance TrustService used when hosting the service in-process. This is
// set up by |CreateInProcessTrustServiceOnThread()| and destroyed by
// |ShutDownTrustService()|.
trust::TrustService* g_in_process_trust_service_instance = nullptr;

void CreateInProcessTrustServiceOnThread(
    mojo::PendingReceiver<trust::mojom::TrustService> receiver) {
#if BUILDFLAG(IS_ANDROID)
  if (base::GetFieldTrialParamByFeatureAsBool(
          features::kBigLittleScheduling,
          features::kBigLittleSchedulingNetworkMainBigParam, false)) {
    SetCpuAffinityForCurrentThread(base::CpuAffinityMode::kBigCoresOnly);
  }
#endif

  g_in_process_trust_service_instance = new trust::TrustService(std::move(receiver));
}

void CreateInProcessTrustService(
    mojo::PendingReceiver<trust::mojom::TrustService> receiver) {
  TRACE_EVENT0("loading", "CreateInProcessTrustService");
  scoped_refptr<base::SingleThreadTaskRunner> task_runner;
 
  base::Thread::Options options(base::MessagePumpType::IO, 0);
  GetTrustServiceDedicatedThread().StartWithOptions(std::move(options));
  task_runner = GetTrustServiceDedicatedThread().task_runner();
  
  GetTrustTaskRunner()->PostTask(
      FROM_HERE, base::BindOnce(&CreateInProcessTrustServiceOnThread,
                                std::move(receiver)));
}

base::RepeatingClosureList& GetCrashHandlersList() {
  static base::NoDestructor<base::RepeatingClosureList> s_list;
  return *s_list;
}

void OnTrustServiceCrash() {
  DCHECK(::content::BrowserThread::CurrentlyOn(::content::BrowserThread::UI));
  DCHECK(g_trust_service_remote);
  DCHECK(g_trust_service_remote->is_bound());
  DCHECK(!g_trust_service_remote->is_connected());
  g_last_trust_service_crash = base::Time::Now();
  GetCrashHandlersList().Notify();
}

}  // namespace

// we always run trust service OOP for now
bool IsInProcessTrustService() {
  return false;
}

trust::mojom::TrustService* GetTrustService() {
  static bool once = false;

  if (!g_trust_service_remote)
    g_trust_service_remote = new mojo::Remote<trust::mojom::TrustService>;

  if (!g_trust_service_remote->is_bound() ||
      !g_trust_service_remote->is_connected()) {

    once = false;    
    bool service_was_bound = g_trust_service_remote->is_bound();
    g_trust_service_remote->reset();

    if (::content::GetContentClient()->browser()->IsShuttingDown()) {
      // This happens at system shutdown, since in other scenarios the trust
      // process would only be torn down once the message loop stopped running.
      // We don't want to start the trust service again so just create message
      // pipe that's not bound to stop consumers from requesting creation of the
      // service.
      auto receiver = g_trust_service_remote->BindNewPipeAndPassReceiver();
      auto leaked_pipe = receiver.PassPipe().release();
    } else {
      if (!g_force_create_trust_service_directly) {
        mojo::PendingReceiver<trust::mojom::TrustService> receiver =
            g_trust_service_remote->BindNewPipeAndPassReceiver();

        g_trust_service_remote->set_disconnect_handler(
            base::BindOnce(&OnTrustServiceCrash));

        if (IsInProcessTrustService()) {
          CreateInProcessTrustService(std::move(receiver));
        } else {
          if (service_was_bound)
            LOG(ERROR) << "Trust service crashed, restarting service.";

          ::content::ServiceProcessHost::Launch(std::move(receiver),
                                     ::content::ServiceProcessHost::Options()
                                         .WithDisplayName(u"Trust Service")
                                         .Pass());
        }
      }
    }
  }

  if (!once) {
    g_trust_service_remote->get()->Launch();
    once = true;
  }

  return g_trust_service_remote->get();
}

base::CallbackListSubscription RegisterTrustServiceCrashHandler(
    base::RepeatingClosure handler) {
  DCHECK(::content::BrowserThread::CurrentlyOn(::content::BrowserThread::UI));
  DCHECK(!handler.is_null());

  return GetCrashHandlersList().Add(std::move(handler));
}

scoped_refptr<base::SequencedTaskRunner>& GetTrustTaskRunnerStorage() {
  static base::NoDestructor<scoped_refptr<base::SequencedTaskRunner>> storage;
  return *storage;
}

const scoped_refptr<base::SequencedTaskRunner>& GetTrustTaskRunner() {
  DCHECK(IsInProcessTrustService());
  return GetTrustTaskRunnerStorage();
}

void ShutDownTrustService() {
  delete g_trust_service_remote;
  g_trust_service_remote = nullptr;
  if (g_in_process_trust_service_instance) {
    GetTrustTaskRunner()->DeleteSoon(FROM_HERE, g_in_process_trust_service_instance);
    g_in_process_trust_service_instance = nullptr;
  }
  GetTrustTaskRunnerStorage().reset();
}

}  // namespace content
} // namespace beacon
