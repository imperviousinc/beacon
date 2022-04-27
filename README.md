<img width="350" src="https://user-images.githubusercontent.com/41967894/164736581-db3d215c-70d6-4ee3-94ba-21fc8a0b989e.svg#gh-light-mode-only" alt="Beacon browser">
<img width="350" src="https://user-images.githubusercontent.com/41967894/164736898-bd00ea5a-b97c-4363-b688-59823622d626.svg#gh-dark-mode-only" alt="Beacon browser">

-------

Note: ⚠️ Beacon is still in beta use at your own risk.

A first-class browsing experience for a decentralized internet built with web technologies and secured without third parties. Trustless HTTPS with native DANE support and a DNSSEC chain secured by a peer-to-peer light client.

<kbd>
<img border="1" width="400" src="https://user-images.githubusercontent.com/41967894/164748866-649c78c7-cd76-4613-9d17-82d382320b98.PNG">
</kbd>

## How it works

- Beacon syncs block headers to retrive a verifiable merkle tree root.
- Requests proofs from peers to retrive a DNSSEC signed zone.
- Performs in-browser DNSSEC validation.
- Verifies certificates with [DANE](https://datatracker.ietf.org/doc/html/rfc6698).

### TODOs

There are still lots of things we'd like to do. Contributions are welcome!

* Android & linux support
* Automatic updates using [Omaha 4](https://docs.google.com/document/d/1VlozzSjriRD5Yn9cLzjTrSvXPkxtq47mk2JejkczAss/edit)  
* Signed binaries for windows
* Widevine support
* DNSSEC prefetching to reduce latency
* DANE support for ICANN domains
* Experiment with embedding a DNSSEC chain in x509 certificates 
 and/or a TLS extension (RFC9102). 
* Experiment with embedding HNS proofs in x509 certificates. 
* Block internal Chrome telemetry & other privacy enhancements
* More tests


Development
-------

This repository does not contain the actual Chromium code it will be fetched using `butil`.

### Get started

Install chromium build dependencies for the target platform and then install `butil`. 

```
$ go install github.com/imperviousinc/beacon/tools/src/butil@latest
```
`butil` is beacon's development utility. It helps you apply patches and do various overrides to chromium. Make sure it's in your path.

```
$ mkdir beacon && cd beacon
$ butil clone
$ butil init
```
This may take a while since `init` will fetch chromium. Once it's done, this repo will be at `src/beacon`


#### Building

```
$ butil build debug
```

#### Updating `butil`

`butil` is just a wrapper around the actual tool. You can make changes to `tools/src/realbutil`
and it will get rebuilt automatically. 


#### Making changes to Chromium

Make your modifications to chromium and when you are ready to transfer those into patches:
Note: This will remove any patches that are no longer in chromium.

```
$ butil patches update
```

To remove a patch just undo the changes in chromium repo and call patches update again.


## Credits

Beacon ports patches from Brave mainly for branding and shares a similar patching format/tooling with [brave-core](https://github.com/brave/brave-core.git)
