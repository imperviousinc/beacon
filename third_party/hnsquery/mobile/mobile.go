package mobile

import (
	"context"
	"errors"
	"fmt"
	lru "github.com/hashicorp/golang-lru"
	"github.com/imperviousinc/hnsquery"
	"github.com/imperviousinc/hnsquery/hip5"
	_ "golang.org/x/mobile/bind"
	"log"
	"os"
	"path"
	"time"
)

const (
	HNSNotReady      int = 0
	HNSNotSynced         = 1
	HNSNoPeers           = 2
	HNSInsecure          = 3
	HNSBogus             = 4
	HNSLookupTimeout     = 5
	HNSPeerTimeout       = 6
	HNSSecure            = 7
	HNSErr               = 8
	HNSOk                = 9
	HNSErrCert           = 10
	HNSErrorUnknown      = -1
)

type HNS struct {
	dataDir    string
	client     *hnsquery.Client
	resolver   *hnsquery.Resolver
	certVerify *hnsquery.DNSCertVerifier

	eth           hip5.Ethereum
	secureChannel bool

	tldMemCache  *lru.Cache
	tldDiskCache *diskCache
	// for tests
	disableNameChecks bool
}

type logWriter struct {
}

func (writer logWriter) Write(bytes []byte) (int, error) {
	return fmt.Print("hns: mobile: " + string(bytes))
}

func init() {
	log.SetFlags(0)
	log.SetOutput(new(logWriter))
}

func NewVerifier(dohURL string) (h *HNS, err error) {
	h = &HNS{}
	if h.dataDir, err = os.UserCacheDir(); err != nil {
		return
	}

	h.dataDir = path.Join(h.dataDir, "Impervious")
	h.secureChannel = true

	if err = os.MkdirAll(h.dataDir, os.ModePerm); err != nil {
		log.Printf("cannot create cache dir: %v", err)
		return
	}

	if h.client, err = hnsquery.NewClient(&hnsquery.Config{
		DataDir: h.dataDir,
	}); err != nil {
		return
	}

	if h.resolver, err = hnsquery.NewResolver(&hnsquery.ResolverConfig{
		Forward: dohURL,
	}); err != nil {
		return
	}

	if h.tldMemCache, err = lru.New(100); err != nil {
		return
	}

	tldCacheDir := path.Join(h.dataDir, "TLDCache")
	if err = os.MkdirAll(tldCacheDir, os.ModePerm); err != nil {
		log.Printf("cannot create tld cache dir: %v", err)
		return
	}

	if h.tldDiskCache, err = newDiskCache(tldCacheDir); err != nil {
		return
	}

	h.resolver.TrustAnchorFunc = getPowTrustAnchor(h)
	if h.certVerify, err = hnsquery.NewDNSCertVerifier(h.resolver); err != nil {
		return
	}

	return
}

func (h *HNS) LaunchTA() error {
	return h.client.Run()
}

func (h *HNS) ShutdownTA() {
	_ = h.client.Destroy()
}

func (h *HNS) Ready() bool {
	return h.client.Ready()
}

func (h *HNS) Progress() float32 {
	return h.client.Progress()
}

func (h *HNS) Height() int {
	return int(h.client.Height())
}

func (h *HNS) PeerCount() int {
	return h.client.PeerCount()
}

func (h *HNS) ActivePeerCount() int {
	return h.client.ActivePeerCount()
}

type VerifyResult struct {
	message string
	code    int
}

func (vr *VerifyResult) Message() string {
	return vr.message
}

func (vr *VerifyResult) Code() int {
	return vr.code
}

func verifyResult(code int, err error) *VerifyResult {
	var msg string
	if err != nil {
		msg = err.Error()
	}

	return &VerifyResult{
		message: msg,
		code:    code,
	}
}

func (h *HNS) VerifyCerts(leaf []byte, port, proto, hostname string) *VerifyResult {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	info := &hnsquery.CertVerifyInfo{
		Host:     hostname,
		Port:     port,
		Protocol: proto,
		RawCerts: [][]byte{leaf},
	}

	if h.disableNameChecks {
		info.DisableNameCheck = true
	}

	ok, err := h.certVerify.Verify(ctx, info)
	if err != nil {
		log.Printf(err.Error())
		switch {
		case errors.Is(err, hnsquery.ErrTimeout):
			return verifyResult(HNSPeerTimeout, err)
		case errors.Is(err, hnsquery.ErrCancelled):
			return verifyResult(HNSLookupTimeout, err)
		case errors.Is(err, hnsquery.ErrNotSynced):
			return verifyResult(HNSNotSynced, err)
		case errors.Is(err, hnsquery.ErrNoPeers):
			return verifyResult(HNSNoPeers, err)
		case errors.Is(err, hnsquery.ErrDNSFatal):
			return verifyResult(HNSBogus, err)
		case errors.Is(err, hnsquery.ErrDNSSECFailed):
			return verifyResult(HNSBogus, err)
		case errors.Is(err, hnsquery.ErrCertVerifyFailed):
			return verifyResult(HNSBogus, err)
		}

		return verifyResult(HNSErrorUnknown, err)
	}

	if !ok {
		return verifyResult(HNSInsecure, nil)
	}

	return verifyResult(HNSSecure, nil)
}
