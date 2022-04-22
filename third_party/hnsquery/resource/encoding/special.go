// Adapted from https://github.com/mslipper/handshake

package encoding

import (
	"errors"
	"io"
	"net"
)

func WriteIP4(w io.Writer, ip net.IP) error {
	data := ip.To4()
	if data == nil {
		return errors.New("invalid IP")
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}

func ReadIP4(r io.Reader) (net.IP, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func WriteIP6(w io.Writer, ip net.IP) error {
	data := ip.To16()
	if data == nil {
		return errors.New("invalid IP")
	}
	if _, err := w.Write(data); err != nil {
		return err
	}
	return nil
}

func ReadIP6(r io.Reader) (net.IP, error) {
	buf := make([]byte, 16)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}
