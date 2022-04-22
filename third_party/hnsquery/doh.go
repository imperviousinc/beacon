package hnsquery

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"time"
)

type dohConn struct {
	endpoint *url.URL
	http     *http.Client

	body *io.Reader
	ctx  context.Context

	deadline time.Time
}

func (d *dohConn) Read(b []byte) (n int, err error) {
	if d.body == nil {
		return 0, io.ErrClosedPipe
	}

	return (*d.body).Read(b)
}

func (d *dohConn) Write(b []byte) (n int, err error) {
	if d.body != nil {
		return 0, io.ErrClosedPipe
	}

	if len(b) < 2 {
		return 0, fmt.Errorf("bad message")
	}

	// trim length
	b = b[2:]
	ctx := d.ctx

	if !d.deadline.IsZero() {
		var cancel context.CancelFunc
		ctx, cancel = context.WithDeadline(d.ctx, d.deadline)
		defer cancel()
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, d.endpoint.String(), bytes.NewBuffer(b))
	if err != nil {
		return 0, fmt.Errorf("failed making http request: %v", err)
	}

	req.Header.Add("Content-Type", "application/dns-message")
	req.Host = d.endpoint.Host

	resp, err := d.http.Do(req)
	if err != nil {
		return 0, fmt.Errorf("failed reading http response: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed with http status code %d", resp.StatusCode)
	}

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("failed reading response: %v", err)
	}

	msg := make([]byte, 2+len(buf))
	binary.BigEndian.PutUint16(msg, uint16(len(buf)))
	copy(msg[2:], buf)

	reader := io.Reader(bytes.NewReader(msg))
	d.body = &reader
	return len(b), nil
}

func (d *dohConn) Close() error {
	return nil
}

func (d *dohConn) LocalAddr() net.Addr {
	return nil
}

func (d *dohConn) RemoteAddr() net.Addr {
	return nil
}

func (d *dohConn) SetDeadline(t time.Time) error {
	d.deadline = t
	return nil
}

func (d *dohConn) SetReadDeadline(t time.Time) error {
	d.deadline = t
	return nil

}

func (d *dohConn) SetWriteDeadline(t time.Time) error {
	d.deadline = t
	return nil
}
