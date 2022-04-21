package internal

import (
	"encoding/hex"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"runtime"

	lru "github.com/hashicorp/golang-lru"
	"github.com/imperviousinc/beacon/components/core/internal/content"
	"github.com/imperviousinc/beacon/components/core/public/proto"
	"github.com/imperviousinc/hnsquery"
	"google.golang.org/grpc"
)

const (
	// TODO: replace with mojo
	trustServiceEndpoint = "44961"
	resourcesEndpoint    = "44962"
)

type Config struct {
	hsq      *hnsquery.Client
	verifier *hnsquery.DNSCertVerifier
	server   *grpc.Server
}

func NewAPI() (*Config, error) {
	var err error
	c := &Config{}

	// create hsq client which is a libhsk binding
	if c.hsq, err = NewHNSQueryClient(); err != nil {
		return nil, err
	}

	// create a cert verifier which is a stub dnssec validating
	// resolver that uses hsq as a trust anchor
	if c.verifier, err = NewCertVerifier("https://hs.dnssec.dev/dns-query", c.hsq); err != nil {
		return nil, err
	}

	c.server = NewGRPCCertVerifierServer(c)
	return c, nil
}

func (c *Config) Launch() {
	hsqLaunch := func() {
		err := c.hsq.Run()
		if err != nil {
			panic(err)
		}
	}

	go hsqLaunch()
	pages := NewContentPages(c)
	go pages.Serve()

	// TODO: replace with sockets
	listen, err := net.Listen("tcp", "127.0.0.1:"+trustServiceEndpoint)
	if err != nil {
		panic(err)
	}

	if err := c.server.Serve(listen); err != nil {
		panic(err)
	}
}

func NewHNSQueryClient() (*hnsquery.Client, error) {
	cacheDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed getting cache dir: %v", err)
	}

	if runtime.GOOS == "windows" {
		cacheDir = filepath.Join(cacheDir, "..", "Local", "Impervious", "Beacon", "ServiceCache")
	} else {
		cacheDir = filepath.Join(cacheDir, "Impervious", "Beacon", "ServiceCache")
	}

	if err = os.MkdirAll(cacheDir, 0700); err != nil {
		return nil, fmt.Errorf("failed making cache dir `%s`: %v", cacheDir, err)
	}

	client, err := hnsquery.NewClient(&hnsquery.Config{DataDir: cacheDir})
	if err != nil {
		return nil, fmt.Errorf("failed creating new hnsquery instance: %v", err)
	}

	return client, nil
}

func NewCertVerifier(dohURL string, q ZoneQuery) (*hnsquery.DNSCertVerifier, error) {
	h := &RootZoneConfig{}
	h.client = q

	var resolver *hnsquery.Resolver
	var err error
	if resolver, err = hnsquery.NewResolver(&hnsquery.ResolverConfig{
		Forward: dohURL,
	}); err != nil {
		return nil, err
	}

	if h.tldMemCache, err = lru.New(100); err != nil {
		return nil, err
	}

	var certVerify *hnsquery.DNSCertVerifier
	resolver.TrustAnchorPointHandler = getPowTrustAnchor(h)
	if certVerify, err = hnsquery.NewDNSCertVerifier(resolver); err != nil {
		return nil, err
	}

	return certVerify, nil
}

func NewGRPCCertVerifierServer(c *Config) *grpc.Server {
	s := grpc.NewServer()
	proto.RegisterCertVerifierServer(s, &CertVerifierGRPC{config: c})
	return s
}

func NewContentPages(c *Config) *content.Config {
	return content.NewContent(resourcesEndpoint, func() *content.HandshakeStatus {
		root := c.hsq.NameRoot()
		return &content.HandshakeStatus{
			TotalPeers:  c.hsq.PeerCount(),
			ActivePeers: c.hsq.ActivePeerCount(),
			Height:      c.hsq.Height(),
			Urkel:       hex.EncodeToString(root[:]),
			Synced:      c.hsq.Ready(),
			Progress:    int(c.hsq.Progress() * 100),
		}
	})
}
