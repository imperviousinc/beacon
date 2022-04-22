package hnsquery

/*
#include "hsk.h"
#include "hns.h"
#include "pool.h"
#include "uv.h"
#include <stdio.h>
#include <time.h>
#include <stdlib.h>

#cgo CFLAGS: -DCGO_BUILD -I${SRCDIR}/build/include -I${SRCDIR}/build/include/hsk -I${SRCDIR}/build/include/hsk/chacha20 -I${SRCDIR}/build/include/hsk/poly1305 -I${SRCDIR}/build/include/hsk/secp256k1
#cgo linux LDFLAGS: ${SRCDIR}/build/lib/libuv.a ${SRCDIR}/build/lib/libhsk.a -ldl
#cgo windows LDFLAGS: ${SRCDIR}/build/lib/libuv.a ${SRCDIR}/build/lib/libhsk.a -lws2_32 -liphlpapi -luserenv -lbcrypt -lpsapi
#cgo ios LDFLAGS: -L${SRCDIR}/build/ios/lib -luv -lhsk
#cgo darwin,!ios LDFLAGS: -L${SRCDIR}/build/lib -luv -lhsk
*/
import "C"
import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"path"
	"sync"
	"time"
	"unsafe"

	"github.com/miekg/dns"
)

var ErrNotSynced = fmt.Errorf("client is still syncing")
var ErrTimeout = fmt.Errorf("request timed out")
var ErrCancelled = fmt.Errorf("operation cancelled")
var ErrNoPeers = fmt.Errorf("no peers")

type Config struct {
	DataDir string
}

type Client struct {
	config *Config
	ctx    *C.hns_ctx
	ctxId  uint64

	callbacks *cgoHSKAccess

	// Shutdown handling
	sync.RWMutex
	started     bool
	closing     chan struct{}
	closed      chan struct{}
	closingOnce sync.Once
}

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print("hns: " + string(bytes))
}

func init() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}

func NewClient(config *Config) (*Client, error) {
	ctx := C.hns_ctx_create()
	if ctx == nil {
		return nil, fmt.Errorf("failed creating context")
	}

	ctxId := rand.Uint64()
	if config.DataDir != "" {
		hdrFile := path.Join(config.DataDir, "chain.bin")
		chdrFile := C.CString(hdrFile)

		C.hns_ctx_set_headers_file(ctx, chdrFile)
		C.free(unsafe.Pointer(chdrFile))
	}

	C.hns_ctx_set_id(ctx, C.ulonglong(ctxId))

	c := &Client{
		config:    config,
		ctx:       ctx,
		ctxId:     ctxId,
		closing:   make(chan struct{}, 1),
		closed:    make(chan struct{}, 1),
		callbacks: newCGOHSK(),
	}

	ctxMap.Lock()
	ctxMap.contexts[ctxId] = c.callbacks
	ctxMap.Unlock()

	return c, nil
}

func getContextId(v unsafe.Pointer) uint64 {
	return uint64(C.hns_ctx_get_id((*C.struct_hns_ctx)(v)))
}

func hskCodeToError(code C.int) error {
	switch code {
	case C.HNS_ETIMEOUT:
		return ErrTimeout
	case C.HNS_ENOPEERS:
		return ErrNoPeers
	case C.HNS_ENOTSYNCED:
		return ErrNotSynced
	case C.HNS_ENOMEM:
		return fmt.Errorf("out of memory")
	}

	return fmt.Errorf("hns error (code: %d)", int(code))
}

func (client *Client) didStart() bool {
	client.RLock()
	defer client.RUnlock()

	return client.started
}

func (client *Client) Run() error {
	if client.didStart() {
		return fmt.Errorf("client already started or no longer valid")
	}

	client.Lock()
	client.started = true
	client.Unlock()

	// starts event loop
	r := C.hns_ctx_start(client.ctx)
	C.hns_ctx_destroy(client.ctx)
	client.ctx = nil

	defer close(client.closed)

	if r != C.HNS_SUCCESS {
		return hskCodeToError(r)
	}

	return nil
}

func (client *Client) Start(ready chan error) {
	stop := make(chan error, 1)

	go func() {
		err := client.Run()
		stop <- err
	}()

	readyTicker := time.NewTicker(300 * time.Millisecond)
	go func() {
		for {
			select {
			case err := <-stop:
				ready <- err
				return
			case <-readyTicker.C:
				if client.Ready() && client.ActivePeerCount() > 0 {
					ready <- nil
					return
				}
			case <-client.closing:
				log.Println("shutting down ticker")
				return
			}
		}
	}()
}

func (client *Client) Destroy() error {
	// only wait for closed event
	// if the client was started
	if client.didStart() {
		client.shutdown()
		<-client.closed
	}

	return nil
}

func (client *Client) shutdown() error {
	// sync once is used here just in case
	// shutdown is called multiple times
	client.closingOnce.Do(func() {
		// any listening go routines
		// should quit
		close(client.closing)

		// signal libuv event loop to stop
		C.hns_ctx_shutdown(client.ctx)

		// clear callbacks
		ctxMap.Lock()
		delete(ctxMap.contexts, client.ctxId)
		ctxMap.Unlock()
	})

	return nil
}

func (client *Client) GetZone(ctx context.Context, name string) (rrs []dns.RR, err error) {
	resultReady := make(chan struct{}, 1)

	var f CallbackFunc = func(res []dns.RR, resErr error) {
		rrs = res
		if resErr != nil {
			err = fmt.Errorf("failed resolving zone %s: %w", name, resErr)
		}

		resultReady <- struct{}{}
	}
	defer client.callbacks.removeCallback(name, &f)

	if doLookup := client.callbacks.addCallback(name, &f); doLookup {
		cname := C.CString(name)
		C.hns_resolve(client.ctx, cname)
		C.free(unsafe.Pointer(cname))
	}

	select {
	case <-ctx.Done():
		err = fmt.Errorf("failed resolving zone %s: %w", name, ErrCancelled)
		return
	case <-client.closed:
		err = fmt.Errorf("client is not running")
		return
	case <-client.closing:
		err = fmt.Errorf("stopped resolving shutting down: %w", ErrCancelled)
		return
	case <-resultReady:
		return
	}
}

func (client *Client) Ready() bool {
	return bool(C.hns_chain_ready(client.ctx))
}

func (client *Client) Progress() float32 {
	return float32(C.hns_chain_progress(client.ctx))
}

func (client *Client) Height() uint64 {
	return uint64(C.hns_chain_height(client.ctx))
}

func (client *Client) PeerCount() int {
	return int(C.hns_pool_total_peers(client.ctx))
}

func (client *Client) ActivePeerCount() int {
	return int(C.hns_pool_active_peers(client.ctx))
}

func (client *Client) NameRoot() []byte {
	root := C.hns_chain_name_root(client.ctx)
	defer C.free(unsafe.Pointer(root))

	return C.GoBytes(unsafe.Pointer(root), 32)
}
