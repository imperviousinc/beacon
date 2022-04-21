#ifndef BEACON_COMPONENTS_CORE_BINDINGS_CORE_LIBRARY_H_
#define BEACON_COMPONENTS_CORE_BINDINGS_CORE_LIBRARY_H_

#include <memory>
#include "base/native_library.h"

namespace beacon {
namespace core {

class CoreLibrary {
  public:
    // Loads the Go library and relevant functions required. 
    // Will return nullptr if it fails.
    static std::unique_ptr<CoreLibrary> Create();

    CoreLibrary(const CoreLibrary&) = delete;
    CoreLibrary& operator=(const CoreLibrary&) = delete;
    ~CoreLibrary();

    // Launches beacon service (blocking) will return an error
    // in case of failure.
    int32_t Launch();


    // Returns whether this instance is valid (i.e. all necessary functions have
    // been loaded.)
    bool IsValid() const;
    void Shutdown();

    using LaunchFunc = int32_t (*)();
    LaunchFunc GetLaunchFunc();
 private:
  CoreLibrary(base::NativeLibrary native_library);

  // Loads the functions exposed by the native library.
  void LoadFunctions();
  base::NativeLibrary native_library_;
  LaunchFunc launch_func_ = nullptr;

  using ShutdownFunc = void (*)();
  ShutdownFunc shutdown_func_ = nullptr;

};

} // namespace core
} // namespace beacon

#endif // BEACON_COMPONENTS_CORE_BINDINGS_CORE_LIBRARY_H_
