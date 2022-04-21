#ifndef BEACON_BROWSER_BEACON_BEACON_PROFILE_PREFS_H_
#define BEACON_BROWSER_BEACON_BEACON_PROFILE_PREFS_H_

namespace user_prefs {
class PrefRegistrySyncable;
}

namespace beacon {

void RegisterProfilePrefs(user_prefs::PrefRegistrySyncable* registry);

}  // namespace beacon

#endif  // BEACON_BROWSER_BEACON_BEACON_PROFILE_PREFS_H_
