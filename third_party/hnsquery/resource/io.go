// Adapted from https://github.com/mslipper/handshake

package resource

import (
	"errors"
	"github.com/miekg/dns"
	"io"
)

type ResourceReader struct {
	b   []byte
	off int
}

func NewResourceReader(b []byte) *ResourceReader {
	return &ResourceReader{
		b: b,
	}
}

func (r *ResourceReader) Read(p []byte) (n int, err error) {
	bufLen := len(r.b)
	if r.off == bufLen {
		return 0, io.EOF
	}

	if r.off+len(p) > bufLen {
		amtRead := bufLen - r.off
		copy(p, r.b[r.off:])
		r.off = bufLen
		return amtRead, nil
	}

	copy(p, r.b[r.off:])
	r.off += len(p)
	return len(p), nil
}

type ResourceWriter struct {
	b   []byte
	off int
}

func NewResourceWriter() *ResourceWriter {
	return &ResourceWriter{
		b: make([]byte, 1024, 1024),
	}
}

func (w *ResourceWriter) Write(p []byte) (int, error) {
	bufLen := len(w.b)
	if w.off == bufLen {
		return 0, errors.New("resource too long")
	}

	if w.off+len(p) > bufLen {
		amtWritten := bufLen - w.off
		copy(w.b[w.off:], p[:amtWritten])
		w.off = bufLen
		return amtWritten, nil
	}

	copy(w.b[w.off:], p)
	w.off += len(p)
	return len(p), nil
}

func (w *ResourceWriter) Bytes() []byte {
	return w.b[:w.off]
}

func readName(r *ResourceReader) (string, error) {
	name, off, err := dns.UnpackDomainName(r.b, r.off)
	if err != nil {
		return "", err
	}
	r.off = off
	return name, nil
}

func writeName(w *ResourceWriter, name string, compressMap map[string]int) error {
	newOff, err := dns.PackDomainName(name, w.b, w.off, compressMap, true)
	if err != nil {
		return err
	}
	w.off = newOff
	return nil
}
