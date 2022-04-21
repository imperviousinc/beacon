#include "beacon/components/core/bindings/core_library.h"
#include "base/logging.h"
#include "base/base_paths.h"
#include "base/path_service.h"
#include "build/build_config.h"
#include "base/compiler_specific.h"
#include "base/memory/ptr_util.h"

#if BUILDFLAG(IS_MAC)
#include "base/mac/bundle_locations.h"
#include "base/mac/foundation_util.h"
#endif

namespace beacon {
namespace core {

// static
std::unique_ptr<CoreLibrary> CoreLibrary::Create() {
  base::FilePath base_dir;
#if BUILDFLAG(IS_MAC)
  if (base::mac::AmIBundled()) {
    base_dir = base::mac::FrameworkBundlePath().Append("Libraries");
  } else {
#endif
    if (!base::PathService::Get(base::DIR_MODULE, &base_dir)) {
      LOG(ERROR) << "Error getting app dir";
      return nullptr;
    }
#if BUILDFLAG(IS_MAC)
  }
#endif

 auto lib_path = base_dir.AppendASCII(
          base::GetNativeLibraryName("beacon"));

 LOG(INFO) << "Loading beacon core library from path: " << lib_path;

  base::NativeLibraryLoadError error;
  base::NativeLibrary native_library = base::LoadNativeLibrary(
      lib_path, &error);
  if (!native_library) {
    LOG(ERROR) << "Failed to initialize beacon core library: "
               << error.ToString();
    return nullptr;
  }

  std::unique_ptr<CoreLibrary>
      core_native_library =
          base::WrapUnique<CoreLibrary>(
              new CoreLibrary(std::move(native_library)));
  if (core_native_library->IsValid()) {
    return core_native_library;
  }
  LOG(ERROR) << "Could not find all required functions for beacon helper "
                "native library";
  return nullptr;
}

DISABLE_CFI_ICALL
void CoreLibrary::LoadFunctions() {
  launch_func_ =
      reinterpret_cast<LaunchFunc>(
          base::GetFunctionPointerFromNativeLibrary(
              native_library_,
              "BeaconHelper_Launch"));

  shutdown_func_ =
        reinterpret_cast<ShutdownFunc>(
            base::GetFunctionPointerFromNativeLibrary(
                native_library_,
                "BeaconHelper_Shutdown"));
}

DISABLE_CFI_ICALL
bool CoreLibrary::IsValid() const {
  return launch_func_ && shutdown_func_;
}

DISABLE_CFI_ICALL
CoreLibrary::LaunchFunc CoreLibrary::GetLaunchFunc() {
  CHECK(IsValid());
  return launch_func_;
}

CoreLibrary::CoreLibrary(
    base::NativeLibrary native_library)
    : native_library_(std::move(native_library)) {
  LoadFunctions();
}
CoreLibrary::~CoreLibrary() = default;

DISABLE_CFI_ICALL
int32_t CoreLibrary::Launch() {
  CHECK(IsValid());
  return launch_func_();
}

DISABLE_CFI_ICALL
void CoreLibrary::Shutdown() {
  CHECK(IsValid());
  return shutdown_func_();
}

} // namespace core
} // namespace beacon
