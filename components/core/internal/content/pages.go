package content

import (
	"embed"
	"encoding/json"
	"net/http"
)

var (
	//go:embed resources
	resources embed.FS
)

type HandshakeStatus struct {
	TotalPeers  int    `json:"totalPeers"`
	ActivePeers int    `json:"activePeers"`
	Height      uint64 `json:"height"`
	Urkel       string `json:"urkel"`
	Synced      bool   `json:"synced"`
	Progress    int    `json:"progress"`
}

type Config struct {
	GetHandshakeStatus func() *HandshakeStatus
	Port               string
}

func NewContent(port string, handler func() *HandshakeStatus) *Config {
	return &Config{GetHandshakeStatus: handler, Port: port}
}

func (c *Config) Serve() error {
	fs := http.FileServer(http.FS(resources))
	http.Handle("/resources/info.json", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Content-Security-Policy",
			"frame-ancestors chrome://welcome chrome://hns-internals;")

		status := c.GetHandshakeStatus()
		resp, err := json.Marshal(status)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(200)
		w.Write(resp)
	}))

	http.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Security-Policy",
			"frame-ancestors chrome://welcome chrome://hns-internals;")
		fs.ServeHTTP(w, req)
	}))

	return http.ListenAndServe("127.0.0.1:"+c.Port, nil)
}
