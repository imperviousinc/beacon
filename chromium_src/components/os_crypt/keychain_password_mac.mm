#include "components/os_crypt/keychain_password_mac.h"

#include <utility>

#include "base/command_line.h"

namespace {

const char kBeaconDefaultServiceName[] = "Beacon Safe Storage";
const char kBeaconDefaultAccountName[] = "Beacon";

KeychainPassword::KeychainNameType& GetBeaconServiceName();
KeychainPassword::KeychainNameType& GetBeaconAccountName();

}

#define BEACON_GET_SERVICE_NAME return GetBeaconServiceName();
#define BEACON_GET_ACCOUNT_NAME return GetBeaconAccountName();
#include "src/components/os_crypt/keychain_password_mac.mm"
#undef BEACON_GET_SERVICE_NAME
#undef BEACON_GET_ACCOUNT_NAME

namespace {

std::pair<std::string, std::string> GetServiceAndAccountName() {
  std::string service_name, account_name;
  base::CommandLine* command_line = base::CommandLine::ForCurrentProcess();
  if (command_line->HasSwitch("import-chrome")) {
    service_name = std::string("Chrome Safe Storage");
    account_name = std::string("Chrome");
  } else if (command_line->HasSwitch("import-chromium") ||
             command_line->HasSwitch("import-beacon")) {
    service_name = std::string("Chromium Safe Storage");
    account_name = std::string("Chromium");
  } else {
    service_name = std::string(kBeaconDefaultServiceName);
    account_name = std::string(kBeaconDefaultAccountName);
  }
  return std::make_pair(service_name, account_name);
}

KeychainPassword::KeychainNameType& GetBeaconServiceName() {
  static KeychainNameContainerType service_name(
      GetServiceAndAccountName().first);
  return *service_name;
}

KeychainPassword::KeychainNameType& GetBeaconAccountName() {
  static KeychainNameContainerType account_name(
      GetServiceAndAccountName().second);
  return *account_name;
}

}
