# Trust Service

This mojo service runs almost always out-of-process possibly even on Android (TODO). It gets started and shutdown along with the Network Service. This service isn't tied to a profile.

It launches HNS peer-to-peer client. It's also responsible for fetching/validating DNSSEC. Depends on //beacon/components/core


### TODOs

* Define a cert verifier factory/cert verifier mojo interfaces and plumb it to the network context.
* DNSSEC Prefetching to reduce latency.
* More testing.
