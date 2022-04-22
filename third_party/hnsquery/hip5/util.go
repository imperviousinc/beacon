package hip5

// some utils adapted from
// https://github.com/ethereum/go-ethereum/blob/release/1.7/contracts/ens/ens.go
// https://github.com/wealdtech/go-ens/blob/904e0feb4c0df8478b11e9e475afde5852c87763/namehash.go

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/miekg/dns"
	"golang.org/x/crypto/sha3"
	"golang.org/x/net/idna"
	"strings"
	"time"
)

var p = idna.New(idna.MapForLookup(), idna.StrictDomainName(false), idna.Transitional(false))

func ensParentNode(name string) (common.Hash, common.Hash) {
	parts := strings.SplitN(name, ".", 2)
	label := crypto.Keccak256Hash([]byte(parts[0]))
	if len(parts) == 1 {
		return [32]byte{}, label
	} else {
		parentNode, parentLabel := ensParentNode(parts[1])
		return crypto.Keccak256Hash(parentNode[:], parentLabel[:]), label
	}
}

func EnsNode(name string) common.Hash {
	parentNode, parentLabel := ensParentNode(name)
	return crypto.Keccak256Hash(parentNode[:], parentLabel[:])
}

// Normalize normalizes a name according to the ENS rules
func Normalize(input string) (output string, err error) {
	output, err = p.ToUnicode(input)
	if err != nil {
		return
	}
	// If the name started with a period then ToUnicode() removes it, but we want to keep it
	if strings.HasPrefix(input, ".") && !strings.HasPrefix(output, ".") {
		output = "." + output
	}
	return
}

// LabelHash generates a simple hash for a piece of a name.
func LabelHash(label string) (hash [32]byte, err error) {
	normalizedLabel, err := Normalize(label)
	if err != nil {
		return
	}

	sha := sha3.NewLegacyKeccak256()
	if _, err = sha.Write([]byte(normalizedLabel)); err != nil {
		return
	}
	sha.Sum(hash[:0])
	return
}

// NameHash generates a hash from a name that can be used to
// look up the name in ENS
func NameHash(name string) (hash [32]byte, err error) {
	if name == "" {
		return
	}
	normalizedName, err := Normalize(name)
	if err != nil {
		return
	}
	parts := strings.Split(normalizedName, ".")
	for i := len(parts) - 1; i >= 0; i-- {
		if hash, err = nameHashPart(hash, parts[i]); err != nil {
			return
		}
	}
	return
}

func nameHashPart(currentHash [32]byte, name string) (hash [32]byte, err error) {
	sha := sha3.NewLegacyKeccak256()
	if _, err = sha.Write(currentHash[:]); err != nil {
		return
	}

	nameSha := sha3.NewLegacyKeccak256()
	if _, err = nameSha.Write([]byte(name)); err != nil {
		return
	}
	nameHash := nameSha.Sum(nil)
	if _, err = sha.Write(nameHash); err != nil {
		return
	}
	sha.Sum(hash[:0])
	return
}

func hashDnsName(name string) ([32]byte, error) {
	var qnameWire [266]byte
	off, err := dns.PackDomainName(name, qnameWire[:], 0, nil, false)
	if err != nil {
		return [32]byte{}, fmt.Errorf("error packing qname `%s`: %v", name, err)
	}

	var qnameHash [32]byte
	hash := crypto.Keccak256(qnameWire[:off])
	copy(qnameHash[:], hash)

	return qnameHash, nil
}

func unpackRRSet(raw []byte) []dns.RR {
	if len(raw) == 0 {
		return nil
	}

	var (
		rrs   []dns.RR
		rr    dns.RR
		err   error
		rrOff = 0
	)

	for {
		rr, rrOff, err = dns.UnpackRR(raw, rrOff)
		if err != nil || rr.Header().Rdlength == 0 {
			break
		}

		rrs = append(rrs, rr)
		if rrOff == len(raw) {
			break
		}
	}

	return rrs
}

func toNode(name string) string {
	return LastNLabels(name, 2)
}

// LastNLabels returns a lower cased string
// with last n labels from the specified domain name
func LastNLabels(name string, n int) string {
	name = dns.CanonicalName(name)
	parts := dns.SplitDomainName(name)
	if len(parts) <= n {
		return strings.Join(parts, ".")
	}

	return strings.Join(parts[len(parts)-n:], ".")
}

// FirstNLabels returns a lower cased string
// with first n labels from the specified domain name
func FirstNLabels(name string, n int) string {
	name = dns.CanonicalName(name)
	parts := dns.SplitDomainName(name)
	if len(parts) <= n {
		return strings.Join(parts, ".")
	}

	return strings.Join(parts[:n], ".")
}

func nsToRR(ns []*dns.NS) (rrs []dns.RR) {
	for _, rr := range ns {
		rrs = append(rrs, rr)
	}
	return
}

// getTTL finds the TTL of an RRSet
// if records have different TTLs it will return the min
// https://datatracker.ietf.org/doc/html/rfc2181#section-5
func getTTL(rrs []dns.RR) time.Duration {
	var ttl uint32 = 10800
	for _, rr := range rrs {
		if ttl > rr.Header().Ttl {
			ttl = rr.Header().Ttl
		}
	}

	// min ttl
	if ttl < 60 {
		return time.Minute
	}

	// max ttl
	if ttl > 10800 {
		return time.Hour * 3
	}

	return time.Duration(ttl) * time.Second
}
