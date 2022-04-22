// Adapted from https://github.com/mslipper/handshake

package encoding

import (
	"encoding/binary"
	"io"
	"math"
)

type Encoder interface {
	Encode(w io.Writer) error
}

func WriteUint64(w io.Writer, val uint64) error {
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, val)
	_, err := w.Write(buf)
	return err
}

func WriteUint32(w io.Writer, val uint32) error {
	buf := make([]byte, 4)
	binary.LittleEndian.PutUint32(buf, val)
	_, err := w.Write(buf)
	return err
}

func WriteUint16(w io.Writer, val uint16) error {
	buf := make([]byte, 2)
	binary.LittleEndian.PutUint16(buf, val)
	_, err := w.Write(buf)
	return err
}

func WriteUInt16BE(w io.Writer, val uint16) error {
	buf := make([]byte, 2)
	binary.BigEndian.PutUint16(buf, val)
	_, err := w.Write(buf)
	return err
}

func WriteUint8(w io.Writer, val uint8) error {
	_, err := w.Write([]byte{val})
	return err
}

func WriteVarint(w io.Writer, val uint64) error {
	var buf []byte
	if val <= 0xfc {
		buf = []byte{uint8(val)}
	} else if val <= math.MaxUint16 {
		buf = make([]byte, 3)
		buf[0] = 0xfd
		binary.LittleEndian.PutUint16(buf[1:], uint16(val))
	} else if val <= math.MaxUint32 {
		buf = make([]byte, 5)
		buf[0] = 0xfe
		binary.LittleEndian.PutUint32(buf[1:], uint32(val))
	} else {
		buf := make([]byte, 9)
		buf[0] = 0xff
		binary.LittleEndian.PutUint64(buf[1:], val)
	}
	_, err := w.Write(buf)
	return err
}

func WriteVarBytes(w io.Writer, buf []byte) error {
	if err := WriteVarint(w, uint64(len(buf))); err != nil {
		return err
	}
	if _, err := w.Write(buf); err != nil {
		return err
	}
	return nil
}
