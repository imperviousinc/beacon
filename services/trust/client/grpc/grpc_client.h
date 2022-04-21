// Copyright 2021 The Chromium Authors. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be
// found in the LICENSE file.
// source: chromeos/services/libassistant/grpc/grpc_libassistant_client.h

#ifndef BEACON_SERVICES_TRUST_CLIENT_GRPC_GRPC_CLIENT_H_
#define BEACON_SERVICES_TRUST_CLIENT_GRPC_GRPC_CLIENT_H_

#include <memory>
#include <string>

#include "base/threading/sequenced_task_runner_handle.h"
#include "beacon/services/trust/client/grpc/grpc_client_thread.h"
#include "beacon/services/trust/client/grpc/grpc_state.h"
#include "beacon/services/trust/client/grpc/grpc_util.h"
// #include "third_party/grpc/src/include/grpcpp/channel.h"
#include <grpcpp/grpcpp.h>
#include "base/logging.h"

namespace beacon {
namespace core {

// Return gRPC method names.
template <typename Request>
std::string GetGrpcMethodName();

// Interface for all methods we as a client can invoke from gRPC
// services. All client methods should be implemented here to send the requests
// to server. We only introduce methods that are currently in use.
class GrpcClient {
 public:
  explicit GrpcClient(std::shared_ptr<grpc::Channel> channel);
  GrpcClient(const GrpcClient&) = delete;
  GrpcClient& operator=(const GrpcClient&) = delete;
  ~GrpcClient();

  // Calls an async client method. ResponseCallback will be invoked from
  // caller's sequence. The raw pointer will be handled by |RPCState| internally
  // and gets deleted upon completion of the RPC call.
  template <typename Request, typename Response>
  void CallServiceMethod(
      const Request& request,
      beacon::core::ResponseCallback<grpc::Status, Response> done,
      beacon::core::StateConfig state_config) {
      LOG(WARNING) << "Calling Service Method Internal";
    
    new beacon::core::RPCState<Response>(
        channel_, client_thread_.completion_queue(),
        GetGrpcMethodName<Request>(), request, std::move(done),
        /*callback_task_runner=*/base::SequencedTaskRunnerHandle::Get(),
        state_config);
  }

 private:
  // This channel will be shared between all stubs used to communicate with
  // multiple services. All channels are reference counted and will be freed
  // automatically.
  std::shared_ptr<grpc::Channel> channel_;

  // Thread running the completion queue.
  beacon::core::GrpcClientThread client_thread_;
};

}  // namespace core
}  // namespace beacon

#endif  // BEACON_SERVICES_TRUST_CLIENT_GRPC_GRPC_CLIENT_H_

