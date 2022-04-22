// Adapted from https://github.com/mslipper/handshake

package encoding

import (
	"encoding/binary"
	"io"
)

type Decoder interface {
	Decode(r io.Reader) error
}

func ReadUint64(r io.Reader) (uint64, error) {
	buf := make([]byte, 8)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(buf), nil
}

func ReadUint32(r io.Reader) (uint32, error) {
	buf := make([]byte, 4)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint32(buf), nil
}

func ReadUint16(r io.Reader) (uint16, error) {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint16(buf), nil
}

func ReadUint16BE(r io.Reader) (uint16, error) {
	buf := make([]byte, 2)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	return binary.BigEndian.Uint16(buf), nil
}

func ReadUint8(r io.Reader) (uint8, error) {
	buf := make([]byte, 1)
	if _, err := io.ReadFull(r, buf); err != nil {
		return 0, err
	}
	return buf[0], nil
}

func ReadVarint(r io.Reader) (uint64, error) {
	sigil, err := ReadByte(r)
	if err != nil {
		return 0, nil
	}
	if sigil < 0xfd {
		return uint64(sigil), nil
	}
	if sigil == 0xfd {
		num := make([]byte, 2)
		if _, err := io.ReadFull(r, num); err != nil {
			return 0, err
		}
		return uint64(binary.LittleEndian.Uint16(num)), nil
	}
	if sigil == 0xfe {
		num := make([]byte, 4)
		if _, err := io.ReadFull(r, num); err != nil {
			return 0, err
		}
		return uint64(binary.LittleEndian.Uint32(num)), nil
	}
	num := make([]byte, 8)
	if _, err := io.ReadFull(r, num); err != nil {
		return 0, err
	}
	return binary.LittleEndian.Uint64(num), nil
}

func ReadBytes(r io.Reader, l int) ([]byte, error) {
	buf := make([]byte, l)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}
	return buf, nil
}

func ReadString(r io.Reader, l int) (string, error) {
	b, err := ReadBytes(r, l)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func ReadByte(r io.Reader) (byte, error) {
	buf, err := ReadBytes(r, 1)
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func ReadVarBytes(r io.Reader) ([]byte, error) {
	l, err := ReadVarint(r)
	if err != nil {
		return nil, err
	}
	return ReadBytes(r, int(l))
}
