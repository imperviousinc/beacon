// Adapted from https://github.com/mslipper/handshake

package resource

import (
	"errors"
	"github.com/imperviousinc/hnsquery/resource/encoding"
	"io"
	"math"
	"net"
)

type RecordType uint8

const (
	RecordTypeDS RecordType = iota
	RecordTypeNS
	RecordTypeGlue4
	RecordTypeGlue6
	RecordTypeSynth4
	RecordTypeSynth6
	RecordTypeTXT
)

func (r RecordType) String() string {
	switch r {
	case RecordTypeDS:
		return "DS"
	case RecordTypeNS:
		return "NS"
	case RecordTypeGlue4:
		return "GLUE4"
	case RecordTypeGlue6:
		return "GLUE6"
	case RecordTypeSynth4:
		return "SYNTH4"
	case RecordTypeSynth6:
		return "SYNTH6"
	case RecordTypeTXT:
		return "TXT"
	default:
		return "unknown"
	}
}

type Record interface {
	Type() RecordType
}

type DSRecord struct {
	KeyTag     uint16
	Algorithm  uint8
	DigestType uint8
	Digest     []byte
}

func (ds *DSRecord) Type() RecordType {
	return RecordTypeDS
}

func (ds *DSRecord) Encode(w io.Writer) error {
	if len(ds.Digest) > math.MaxUint8 {
		return errors.New("digest must be less than 256 bytes")
	}
	if err := encoding.WriteUInt16BE(w, ds.KeyTag); err != nil {
		return err
	}
	if err := encoding.WriteUint8(w, ds.Algorithm); err != nil {
		return err
	}
	if err := encoding.WriteUint8(w, ds.DigestType); err != nil {
		return err
	}
	if err := encoding.WriteUint8(w, uint8(len(ds.Digest))); err != nil {
		return err
	}
	if _, err := w.Write(ds.Digest); err != nil {
		return err
	}
	return nil
}

func (ds *DSRecord) Decode(r io.Reader) error {
	keyTag, err := encoding.ReadUint16BE(r)
	if err != nil {
		return err
	}
	algorithm, err := encoding.ReadUint8(r)
	if err != nil {
		return err
	}
	digestType, err := encoding.ReadUint8(r)
	if err != nil {
		return err
	}
	digestLen, err := encoding.ReadUint8(r)
	if err != nil {
		return err
	}
	digest, err := encoding.ReadBytes(r, int(digestLen))
	if err != nil {
		return err
	}
	ds.KeyTag = keyTag
	ds.Algorithm = algorithm
	ds.DigestType = digestType
	ds.Digest = digest
	return nil
}

type NSRecord struct {
	NS string
}

func (ns *NSRecord) Type() RecordType {
	return RecordTypeNS
}

func (ns *NSRecord) Encode(w *ResourceWriter, compressMap map[string]int) error {
	return writeName(w, ns.NS, compressMap)
}

func (ns *NSRecord) Decode(r *ResourceReader) error {
	server, err := readName(r)
	if err != nil {
		return err
	}
	ns.NS = server
	return nil
}

type Glue4Record struct {
	NS      string
	Address net.IP
}

func (g *Glue4Record) Type() RecordType {
	return RecordTypeGlue4
}

func (g *Glue4Record) Encode(w *ResourceWriter, compressMap map[string]int) error {
	if err := writeName(w, g.NS, compressMap); err != nil {
		return err
	}
	if err := encoding.WriteIP4(w, g.Address); err != nil {
		return err
	}
	return nil
}

func (g *Glue4Record) Decode(r *ResourceReader) error {
	server, err := readName(r)
	if err != nil {
		return err
	}
	addr, err := encoding.ReadIP4(r)
	if err != nil {
		return err
	}
	g.NS = server
	g.Address = addr
	return nil
}

type Glue6Record struct {
	NS      string
	Address net.IP
}

func (g *Glue6Record) Type() RecordType {
	return RecordTypeGlue6
}

func (g *Glue6Record) Encode(w *ResourceWriter, compressMap map[string]int) error {
	if err := writeName(w, g.NS, compressMap); err != nil {
		return err
	}
	if err := encoding.WriteIP6(w, g.Address); err != nil {
		return err
	}
	return nil
}

func (g *Glue6Record) Decode(r *ResourceReader) error {
	server, err := readName(r)
	if err != nil {
		return err
	}
	addr, err := encoding.ReadIP6(r)
	if err != nil {
		return err
	}
	g.NS = server
	g.Address = addr
	return nil
}

type Synth4Record struct {
	Address net.IP
}

func (s *Synth4Record) Type() RecordType {
	return RecordTypeSynth4
}

func (s *Synth4Record) Encode(w io.Writer) error {
	if err := encoding.WriteIP4(w, s.Address); err != nil {
		return err
	}
	return nil
}

func (s *Synth4Record) Decode(r io.Reader) error {
	addr, err := encoding.ReadIP4(r)
	if err != nil {
		return err
	}
	s.Address = addr
	return nil
}

type Synth6Record struct {
	Address net.IP
}

func (s *Synth6Record) Type() RecordType {
	return RecordTypeSynth6
}

func (s *Synth6Record) Encode(w io.Writer) error {
	if err := encoding.WriteIP6(w, s.Address); err != nil {
		return err
	}
	return nil
}

func (s *Synth6Record) Decode(r io.Reader) error {
	addr, err := encoding.ReadIP6(r)
	if err != nil {
		return err
	}
	s.Address = addr
	return nil
}

type TXTRecord struct {
	Entries []string
}

func (t *TXTRecord) Type() RecordType {
	return RecordTypeTXT
}

func (t *TXTRecord) Encode(w io.Writer) error {
	if len(t.Entries) > math.MaxUint8 {
		return errors.New("can encode a max of 255 entries")
	}
	if err := encoding.WriteUint8(w, uint8(len(t.Entries))); err != nil {
		return err
	}
	for _, entry := range t.Entries {
		if len(entry) > math.MaxUint8 {
			return errors.New("entry must be shorter than 256 bytes")
		}
		if err := encoding.WriteUint8(w, uint8(len(entry))); err != nil {
			return err
		}
		if _, err := io.WriteString(w, entry); err != nil {
			return err
		}
	}
	return nil
}

func (t *TXTRecord) Decode(r io.Reader) error {
	count, err := encoding.ReadUint8(r)
	if err != nil {
		return err
	}
	entries := make([]string, count)
	for i := 0; i < int(count); i++ {
		entryLen, err := encoding.ReadUint8(r)
		if err != nil {
			return err
		}
		entry, err := encoding.ReadString(r, int(entryLen))
		if err != nil {
			return err
		}
		entries[i] = entry
	}
	t.Entries = entries
	return nil
}
