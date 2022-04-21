package content

import "testing"

func TestNewContent(t *testing.T) {
	c := NewContent("8290", func() *HandshakeStatus {
		return &HandshakeStatus{
			TotalPeers:  10,
			ActivePeers: 2,
			Height:      1000000,
			Urkel:       "0000 0000",
			Synced:      false,
			Progress:    5,
		}
	})

	c.Serve()
}
