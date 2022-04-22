// Adapted from https://github.com/mslipper/handshake

package resource

import (
	"errors"
	"github.com/imperviousinc/hnsquery/resource/encoding"
	"io"
	"io/ioutil"
)

type CompressorEncoder interface {
	Encode(w *ResourceWriter, compressMap map[string]int) error
	Decode(r *ResourceReader) error
}

type Resource struct {
	TTL     int
	Records []Record
}

func (rs *Resource) Encode(w io.Writer) error {
	rw := NewResourceWriter()
	if err := encoding.WriteUint8(rw, 0); err != nil {
		return err
	}
	compMap := make(map[string]int)
	for _, record := range rs.Records {
		if err := encoding.WriteUint8(rw, uint8(record.Type())); err != nil {
			return err
		}
		switch rt := record.(type) {
		case encoding.Encoder:
			if err := rt.Encode(rw); err != nil {
				return err
			}
		case CompressorEncoder:
			if err := rt.Encode(rw, compMap); err != nil {
				return err
			}
		default:
			return errors.New("cannot encode record")
		}
	}
	if _, err := w.Write(rw.Bytes()); err != nil {
		return err
	}
	return nil
}

func (rs *Resource) Decode(r io.Reader) error {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return err
	}

	rr := NewResourceReader(buf)
	version, err := encoding.ReadUint8(rr)
	if err != nil {
		return err
	}
	if version != 0 {
		return errors.New("invalid serialization version")
	}

	for {
		recType, err := encoding.ReadUint8(rr)
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		var record Record
		switch RecordType(recType) {
		case RecordTypeDS:
			record = new(DSRecord)
		case RecordTypeNS:
			record = new(NSRecord)
		case RecordTypeGlue4:
			record = new(Glue4Record)
		case RecordTypeGlue6:
			record = new(Glue6Record)
		case RecordTypeSynth4:
			record = new(Synth4Record)
		case RecordTypeSynth6:
			record = new(Synth6Record)
		case RecordTypeTXT:
			record = new(TXTRecord)
		default:
			return errors.New("unknown record type")
		}

		switch rt := record.(type) {
		case encoding.Decoder:
			if err := rt.Decode(rr); err != nil {
				return err
			}
		case CompressorEncoder:
			if err := rt.Decode(rr); err != nil {
				return err
			}
		default:
			return errors.New("cannot decode record")
		}
		rs.Records = append(rs.Records, record)
	}
}
