package dnssec

// Modern DNSSEC validation library loosely based on
// https://gitlab.nic.cz/knot/knot-resolver/-/tree/master/lib/dnssec
// https://github.com/semihalev/sdns
//
import (
	"bytes"
	"context"
	"crypto/rsa"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/miekg/dns"
	"math/big"
	"strings"
	"time"
)

var (
	ErrNoDNSKEY               = errors.New("no valid dnskey records found")
	ErrBadDS                  = errors.New("DS record doesn't match zone name")
	ErrNoSignatures           = errors.New("no rrsig records for zone that should be signed")
	ErrMissingDNSKEY          = errors.New("no matching dnskey found for rrsig records")
	ErrSignatureBailiwick     = errors.New("rrsig record out of bailiwick")
	ErrInvalidSignaturePeriod = errors.New("incorrect signature validity period")
	ErrMissingSigned          = errors.New("signed records are missing")
)

// supported dnssec algorithms weaker/unsupported algorithms are treated as unsigned
var supportedAlgorithms = []uint8{dns.RSASHA256, dns.RSASHA512, dns.ECDSAP256SHA256, dns.ECDSAP384SHA384, dns.ED25519}
var supportedDigests = []uint8{dns.SHA256, dns.SHA384}

// DefaultMinRSAKeySize the minimum RSA key size
// that can be used to securely verify messages
const DefaultMinRSAKeySize = 2048

type Zone struct {
	// Name the zone name
	Name string

	// TrustAnchors anchors to validate DNSKeys for this zone
	TrustAnchors []*dns.DS

	// Keys validated DNSKEY RRSet
	Keys map[uint16]*dns.DNSKEY

	// CurrentTime is used to check the validity of all signatures.
	// If zero, the current time is used.
	CurrentTime time.Time

	// The zone expire time
	Expire time.Time

	// MinRSA minimum accepted RSA key size
	MinRSA int

	// For custom zone verification
	VerifyCallback func(ctx context.Context, msg *dns.Msg) (bool, error)
}

func (z *Zone) Secure() bool {
	return z.Keys != nil || z.VerifyCallback != nil
}

func filterDS(zone string, dsSet []dns.RR) ([]*dns.DS, error) {
	if !dns.IsFqdn(zone) {
		return nil, fmt.Errorf("zone must be fqdn")
	}

	type dsKey struct {
		keyTag    uint16
		algorithm uint8
	}

	supported := make(map[dsKey]*dns.DS)
	for _, rr := range dsSet {
		if !strings.EqualFold(zone, rr.Header().Name) {
			return nil, ErrBadDS
		}

		ds, ok := rr.(*dns.DS)
		if !ok {
			continue
		}

		if !isAlgorithmSupported(ds.Algorithm) ||
			!isDigestSupported(ds.DigestType) {
			continue
		}

		key := dsKey{
			keyTag:    ds.KeyTag,
			algorithm: ds.Algorithm,
		}

		// pick strongest supported digest type
		if ds2, ok := supported[key]; ok {
			if ds2.DigestType >= ds.DigestType {
				continue
			}
		}

		supported[key] = ds
	}

	var values []*dns.DS
	for _, rr := range supported {
		values = append(values, rr)
	}

	return values, nil
}

func fromBase64(s []byte) (buf []byte, err error) {
	buflen := base64.StdEncoding.DecodedLen(len(s))
	buf = make([]byte, buflen)
	n, err := base64.StdEncoding.Decode(buf, s)
	buf = buf[:n]
	return
}

func shouldDowngradeKey(k *dns.DNSKEY, minKeySize int) bool {
	if k.Algorithm != dns.RSASHA512 && k.Algorithm != dns.RSASHA256 {
		return false
	}

	// extracted from miekg/dns to check if
	// an exponent is supported by the crypto package
	keybuf, err := fromBase64([]byte(k.PublicKey))
	if err != nil {
		return false
	}

	if len(keybuf) < 1+1+64 {
		// Exponent must be at least 1 byte and modulus at least 64
		return false
	}

	// RFC 2537/3110, section 2. RSA Public KEY Resource Records
	// Length is in the 0th byte, unless its zero, then it
	// it in bytes 1 and 2 and its a 16 bit number
	explen := uint16(keybuf[0])
	keyoff := 1
	if explen == 0 {
		explen = uint16(keybuf[1])<<8 | uint16(keybuf[2])
		keyoff = 3
	}

	if explen > 4 {
		// Exponent larger than supported by the crypto package
		return true
	}

	if explen == 0 || keybuf[keyoff] == 0 {
		// Exponent empty, or contains prohibited leading zero.
		return false
	}

	modoff := keyoff + int(explen)
	modlen := len(keybuf) - modoff
	if modlen < 64 || modlen > 512 || keybuf[modoff] == 0 {
		// Modulus is too small, large, or contains prohibited leading zero.
		return false
	}

	pubkey := new(rsa.PublicKey)

	var expo uint64
	// The exponent of length explen is between keyoff and modoff.
	for _, v := range keybuf[keyoff:modoff] {
		expo <<= 8
		expo |= uint64(v)
	}
	if expo > 1<<31-1 {
		// Larger exponent than supported by the crypto package.
		return true
	}

	pubkey.E = int(expo)
	pubkey.N = new(big.Int).SetBytes(keybuf[modoff:])

	// downgrade if using a weak key size
	if pubkey.N.BitLen() < minKeySize {
		return true
	}

	return false
}

func isAlgorithmSupported(algo uint8) bool {
	for _, curr := range supportedAlgorithms {
		if algo == curr {
			return true
		}
	}

	return false
}

func isDigestSupported(digest uint8) bool {
	for _, curr := range supportedDigests {
		if digest == curr {
			return true
		}
	}

	return false
}

func NewZone(name string, trustAnchors []dns.RR) (*Zone, error) {
	name = dns.CanonicalName(name)

	var dsSet []*dns.DS
	var err error

	if dsSet, err = filterDS(name, trustAnchors); err != nil {
		return nil, err
	}

	return &Zone{
		Name:         name,
		TrustAnchors: dsSet,
		Keys:         nil,
		CurrentTime:  time.Time{},
		MinRSA:       DefaultMinRSAKeySize,
	}, nil
}

// VerifyDNSKeys verifies DNSKEYS from the message using
// the zone trust anchors
func (z *Zone) VerifyDNSKeys(msg *dns.Msg) (map[uint16]*dns.DNSKEY, error) {
	// This zone is insecure
	if len(z.TrustAnchors) == 0 {
		return nil, nil
	}

	matchingKeys := make(map[uint16]*dns.DNSKEY)

	for _, ds := range z.TrustAnchors {
		for _, rr := range msg.Answer {
			if rr.Header().Rrtype != dns.TypeDNSKEY {
				continue
			}

			// simple checks
			key := rr.(*dns.DNSKEY)
			if key.Protocol != 3 {
				continue
			}
			if key.Flags != 256 && key.Flags != 257 {
				continue
			}
			if key.Algorithm != ds.Algorithm {
				continue
			}

			tag := key.KeyTag()
			if tag != ds.KeyTag {
				continue
			}

			dsFromKey := key.ToDS(ds.DigestType)
			if dsFromKey == nil {
				continue
			}

			if !strings.EqualFold(dsFromKey.Digest, ds.Digest) {
				continue
			}

			// we have a valid key
			matchingKeys[tag] = key

		}
	}

	if len(matchingKeys) == 0 {
		return nil, ErrNoDNSKEY
	}

	validKeys := make(map[uint16]*dns.DNSKEY)

	for _, key := range matchingKeys {
		if !shouldDowngradeKey(key, z.MinRSA) {
			validKeys[key.KeyTag()] = key
		}
	}

	if len(validKeys) == 0 {
		return nil, nil
	}

	z.Keys = validKeys
	// verifySignatures will clean up the answer
	// section in the msg with only the valid rr sets
	secure, err := z.verifySignatures(msg, z.Name)
	if err != nil {
		return nil, err
	}

	if !secure {
		return nil, nil
	}

	if len(msg.Answer) == 0 {
		return nil, ErrNoDNSKEY
	}

	trustedKeys := make(map[uint16]*dns.DNSKEY)
	for _, rr := range msg.Answer {
		if rr.Header().Rrtype == dns.TypeDNSKEY {
			key := rr.(*dns.DNSKEY)
			trustedKeys[key.KeyTag()] = key
		}
	}

	return trustedKeys, nil
}

// IsSubDomainStrict checks if child is indeed a child of the parent. If child and parent
// are the same domain false is returned.
func IsSubDomainStrict(parent, child string) bool {
	parentLabels := dns.CountLabel(parent)
	childLabels := dns.CountLabel(child)

	return dns.CompareDomainName(parent, child) == parentLabels &&
		childLabels > parentLabels
}

func (z *Zone) VerifyRRSet(name string, t uint16, set []dns.RR) ([]dns.RR, bool, error) {
	if len(set) == 0 {
		return nil, false, fmt.Errorf("empty set")
	}

	msg := new(dns.Msg)
	msg.SetQuestion(name, t)
	msg.Answer = set
	secure, err := z.verifySignatures(msg, name)
	return msg.Answer, secure, err
}

// verifySignatures verifies signatures in a message
// and removes any invalid rr sets
func (z *Zone) verifySignatures(msg *dns.Msg, qname string) (bool, error) {
	type rrsetId struct {
		owner string
		t     uint16
	}

	downgrade := false
	var lastErr error

	sections := [][]dns.RR{msg.Answer, msg.Ns, msg.Extra}

	// clear sections
	// will fill those as we validate
	msg.Answer = []dns.RR{}
	msg.Ns = []dns.RR{}
	msg.Extra = []dns.RR{}
	var delegations []dns.RR

	for sectionId, section := range sections {
		if len(section) == 0 {
			continue
		}

		verifiedSets := make(map[rrsetId]struct{})

		// Look for all signatures some may be invalid
		// we only need a single valid signature per RRSet
		// this will be used to "discover" covered sets
		// if some records don't have a signature they must be
		// removed from the section
		for _, rr := range section {
			// keep delegations since they are unsigned
			// they will be verified later
			if sectionId == 1 /* authority section */ &&
				rr.Header().Rrtype == dns.TypeNS {
				// must be in bailiwick and within qname
				if IsSubDomainStrict(z.Name, rr.Header().Name) &&
					dns.IsSubDomain(rr.Header().Name, qname) {
					delegations = append(delegations, rr)
				}

				continue
			}

			if rr.Header().Rrtype == dns.TypeRRSIG {
				if sig, ok := rr.(*dns.RRSIG); ok {
					sigName := dns.CanonicalName(sig.Header().Name)

					// if another sig verified the set ignore this one
					if _, ok := verifiedSets[rrsetId{sigName, sig.TypeCovered}]; ok {
						continue
					}

					// we don't care about signatures not in bailiwick
					if !dns.IsSubDomain(z.Name, sigName) {
						lastErr = ErrSignatureBailiwick
						continue
					}

					// look for any valid keys for this signature
					key, ok := z.Keys[sig.KeyTag]
					// RFC4035 5.3.1 bullet 2 signer name must match the name of the zone
					if !ok || !strings.EqualFold(key.Header().Name, sig.SignerName) {
						lastErr = ErrMissingDNSKEY
						continue
					}

					// uses a key that can be downgraded
					// it should fallback to insecure
					// if there are no other secure
					// signatures that can verify the set
					if shouldDowngradeKey(key, z.MinRSA) {
						downgrade = true
						continue
					}

					// extract set covered by signature
					rrset := extractRRSet(section, sig)
					if len(rrset) == 0 {
						lastErr = ErrMissingSigned
						continue
					}

					if err := sig.Verify(key, rrset); err != nil {
						lastErr = err
						continue
					}

					if !sig.ValidityPeriod(z.CurrentTime) {
						lastErr = ErrInvalidSignaturePeriod
						continue
					}

					// verified
					verifiedSets[rrsetId{sigName, sig.TypeCovered}] = struct{}{}

					if sectionId == 0 {
						msg.Answer = append(msg.Answer, rrset...)
						msg.Answer = append(msg.Answer, sig)
						continue
					}

					if sectionId == 1 {
						msg.Ns = append(msg.Ns, rrset...)
						msg.Ns = append(msg.Ns, sig)
						continue
					}

					msg.Extra = append(msg.Extra, rrset...)
					msg.Extra = append(msg.Extra, sig)
				}
			}
		}
	}

	if len(msg.Answer) > 0 || len(msg.Ns) > 0 {
		// append any unsigned delegations
		// to the authority section
		if len(delegations) > 0 {
			msg.Ns = append(msg.Ns, delegations...)
		}

		return true, nil
	}

	// we don't have any secure validation paths
	// if its okay to downgrade mark zone as insecure
	if downgrade {
		msg.Answer = sections[0]
		msg.Ns = sections[1]
		msg.Extra = sections[2]
		return false, nil
	}

	if lastErr != nil {
		return false, fmt.Errorf("error verifying signatures: %v", lastErr)
	}

	return false, ErrNoSignatures
}

func cleanInsecureMsg(msg *dns.Msg) {
	msg.AuthenticatedData = false
	var answer []dns.RR
	var ns []dns.RR
	var extra []dns.RR

	if len(msg.Answer) > 0 {
		msg.Rcode = dns.RcodeSuccess
	}

	for _, rr := range msg.Answer {
		if rr.Header().Rrtype != dns.TypeRRSIG {
			answer = append(answer, rr)
		}
	}

	for _, rr := range msg.Ns {
		if rr.Header().Rrtype != dns.TypeRRSIG {
			ns = append(ns, rr)
		}
	}

	for _, rr := range msg.Extra {
		if rr.Header().Rrtype != dns.TypeRRSIG {
			extra = append(extra, rr)
		}
	}

	msg.Answer = answer
	msg.Ns = ns
	msg.Extra = extra
}

func (z *Zone) Verify(ctx context.Context, msg *dns.Msg, qname string, qtype uint16) (bool, error) {
	if !dns.IsFqdn(z.Name) || !dns.IsFqdn(qname) {
		return false, fmt.Errorf("zone and qname must be fqdn")
	}

	if len(msg.Question) != 1 ||
		!strings.EqualFold(msg.Question[0].Name, qname) ||
		msg.Question[0].Qtype != qtype {
		return false, fmt.Errorf("question mismatch qname %s != %s", qname, msg.Question[0].Name)
	}

	if z.VerifyCallback != nil {
		return z.VerifyCallback(ctx, msg)
	}

	msg.AuthenticatedData = false
	msg.CheckingDisabled = false

	// insecure zone
	if !z.Secure() {
		cleanInsecureMsg(msg)
		return false, nil
	}

	secure, err := z.verifySignatures(msg, qname)
	if err != nil {
		return false, err
	}
	if !secure {
		return false, nil
	}

	// signatures are good verify answer
	if msg.Rcode == dns.RcodeSuccess {
		if len(msg.Answer) == 0 {
			return z.verifyNoData(msg, qname, qtype)
		}

		return verifyAnswer(msg, qname, qtype)
	}

	if msg.Rcode == dns.RcodeNameError {
		return z.verifyNameError(msg, qname)
	}

	return false, fmt.Errorf("unexpected rcode %v", msg.Rcode)
}

// verifyAnswer pass a verified msg with fqdn canonical qname
func verifyAnswer(msg *dns.Msg, qname string, qtype uint16) (bool, error) {
	if len(msg.Answer) == 0 {
		return false, errors.New("empty answer")
	}

	wildcard := false
	nx := false
	labels := uint8(dns.CountLabel(qname))

	// sanitized answer section
	var answer []dns.RR

	for _, rr := range msg.Answer {
		t := rr.Header().Rrtype
		owner := rr.Header().Name

		if t == qtype || t == dns.TypeCNAME {
			// only include rrs that match owner name
			// TODO: flatten CNAMEs if possible
			if strings.EqualFold(qname, owner) {
				answer = append(answer, rr)
			}
			continue
		}

		if t == dns.TypeRRSIG && strings.EqualFold(qname, owner) {
			sig := rr.(*dns.RRSIG)
			if sig.TypeCovered != qtype &&
				sig.TypeCovered != dns.TypeCNAME {
				continue
			}

			answer = append(answer, rr)
			if sig.Labels < labels {
				wildcard = true
			}
			continue
		}
	}

	if len(answer) == 0 {
		return false, errors.New("empty answer")
	}

	msg.Answer = answer

	// if the rrsig is for a wildcard
	// there must be an NSEC proving the original name
	// doesn't exist
	if wildcard {
		for _, rr := range msg.Ns {
			if rr.Header().Rrtype != dns.TypeNSEC {
				continue
			}

			nsec := rr.(*dns.NSEC)
			if nx = covers(nsec.Header().Name, nsec.NextDomain, qname); nx {
				break
			}
		}

		if !nx {
			return false, fmt.Errorf("bad wildcard substitution")
		}
	}

	return true, nil
}

func (z *Zone) verifyNoData(msg *dns.Msg, qname string, qtype uint16) (bool, error) {
	if len(msg.Ns) == 0 {
		return false, fmt.Errorf("no nsec records found")
	}

	for _, rr := range msg.Ns {
		// no authenticated denial of existence
		// for NSEC3 for now it should be downgraded
		if rr.Header().Rrtype == dns.TypeNSEC3 {
			// must be in bailiwick already checked
			// by verifySignatures
			if dns.IsSubDomain(z.Name, rr.Header().Name) {
				return false, nil
			}
		}

		if rr.Header().Rrtype == dns.TypeDS {
			hasNs := false

			if !IsSubDomainStrict(z.Name, rr.Header().Name) {
				return false, fmt.Errorf("ds record must be a child of zone %s", z.Name)
			}

			// NS records aren't signed
			// the owner name must still match
			// the DS record.
			for _, ns := range msg.Ns {
				if ns.Header().Rrtype == dns.TypeNS {
					hasNs = true
					if !strings.EqualFold(ns.Header().Name, rr.Header().Name) {
						return false, fmt.Errorf("bad referral DS owner doesn't match NS")
					}
				}
			}

			// secure delegation with valid
			// NS records
			if hasNs {
				return true, nil
			}

			return false, fmt.Errorf("DS record exists without a delegation")
		}

		if rr.Header().Rrtype == dns.TypeNSEC {
			if nsec, ok := rr.(*dns.NSEC); ok {
				// RFC4035 5.4 bullet 1
				if !strings.EqualFold(nsec.Header().Name, qname) {
					// owner name doesn't match
					// RFC4035 5.4 bullet 2
					return z.verifyNameError(msg, qname)
				}

				// nsec matches qname
				// next domain must be in bailiwick
				if !dns.IsSubDomain(z.Name, nsec.NextDomain) {
					continue
				}

				hasDelegation := false
				hasDS := false

				for _, t := range nsec.TypeBitMap {
					if t == qtype {
						return false, fmt.Errorf("type exists")
					}
					if t == dns.TypeCNAME {
						return false, fmt.Errorf("cname exists")
					}

					if t == dns.TypeDS {
						hasDS = true
						continue
					}

					if t == dns.TypeNS {
						hasDelegation = true
					}
				}

				// verify delegation
				for _, nsRR := range msg.Ns {
					if nsRR.Header().Rrtype == dns.TypeNS {
						if hasDS {
							return false, fmt.Errorf("bad insecure delegation proof " +
								"DS exists in NSEC bitmap")
						}
						if !hasDelegation {
							return false, fmt.Errorf("NS isn't set in NSEC bitmap")
						}
						if !strings.EqualFold(nsRR.Header().Name, nsec.Header().Name) {
							return false, fmt.Errorf("invalid NS owner name")
						}
						if strings.EqualFold(nsRR.Header().Name, z.Name) {
							return false, fmt.Errorf("bad referral")
						}
					}
				}

				return true, nil
			}
		}
	}

	return false, fmt.Errorf("no valid nsec records found")
}

func (z *Zone) verifyNameError(msg *dns.Msg, qname string) (bool, error) {
	nameProof := false
	wildcardProof := false
	qnameParts := dns.SplitDomainName(qname)
	for _, rr := range msg.Ns {
		if nameProof && wildcardProof {
			break
		}

		if rr.Header().Rrtype == dns.TypeNSEC {
			nsec, ok := rr.(*dns.NSEC)
			if !ok {
				continue
			}

			if !nameProof && covers(nsec.Header().Name, nsec.NextDomain, qname) {
				nameProof = true
			}

			if c, err := canonicalNameCompare(nsec.Header().Name, nsec.NextDomain); err == nil && c < 0 {
				if IsSubDomainStrict(qname, nsec.NextDomain) {
					wildcardProof = true
					continue
				}
			}

			if !wildcardProof {
				// find closest wildcard proof that covers qname
				i := 1
				for {
					if len(qnameParts) < i {
						break
					}

					domain := dns.Fqdn("*." + strings.Join(qnameParts[i:], "."))
					if !dns.IsSubDomain(z.Name, domain) {
						break
					}
					if covers(nsec.Header().Name, nsec.NextDomain, domain) {
						wildcardProof = true
						break
					}
					i++
				}
			}
		}
	}

	if !nameProof {
		return false, fmt.Errorf("missing name proof for %s", qname)
	}

	if !wildcardProof {
		return false, fmt.Errorf("missing wildcard proof for %s", qname)
	}

	return true, nil
}

// RFC4034 6.1. Canonical DNS Name Order
// https://tools.ietf.org/html/rfc4034#section-6.1
// Returns -1 if name1 comes before name2, 1 if name1 comes after name2, and 0 if they are equal.
func canonicalNameCompare(name1 string, name2 string) (int, error) {
	// TODO: optimize comparison
	name1 = dns.Fqdn(name1)
	name2 = dns.Fqdn(name2)

	if _, ok := dns.IsDomainName(name1); !ok {
		return 0, errors.New("invalid domain name")
	}
	if _, ok := dns.IsDomainName(name2); !ok {
		return 0, errors.New("invalid domain name")
	}

	labels1 := dns.SplitDomainName(name1)
	labels2 := dns.SplitDomainName(name2)

	var buf1, buf2 [64]byte
	// start comparison from the right
	currentLabel1, currentLabel2, min := len(labels1)-1, len(labels2)-1, 0

	if min = currentLabel1; min > currentLabel2 {
		min = currentLabel2
	}

	for i := min; i > -1; i-- {
		off1, err := dns.PackDomainName(labels1[currentLabel1]+".", buf1[:], 0, nil, false)
		if err != nil {
			return 0, err
		}

		off2, err := dns.PackDomainName(labels2[currentLabel2]+".", buf2[:], 0, nil, false)
		if err != nil {
			return 0, err
		}

		currentLabel1--
		currentLabel2--

		// if the two labels at the same index aren't equal return result
		if res := bytes.Compare(bytes.ToLower(buf1[1:off1-1]),
			bytes.ToLower(buf2[1:off2-1])); res != 0 {
			return res, nil
		}
	}

	// all labels are equal name with least labels is the smallest
	if len(labels1) == len(labels2) {
		return 0, nil
	}

	if len(labels1)-1 == min {
		return -1, nil
	}

	return 1, nil
}

func covers(owner, next, qname string) (result bool) {
	var errs int

	// qname is equal to or before owner can't be covered
	if compareWithErrors(qname, owner, &errs) <= 0 {
		return false
	}

	lastNSEC := compareWithErrors(owner, next, &errs) >= 0
	inRange := lastNSEC || compareWithErrors(qname, next, &errs) < 0
	if !inRange {
		return false
	}

	if errs > 0 {
		return false
	}

	return true
}

func compareWithErrors(a, b string, errs *int) int {
	res, err := canonicalNameCompare(a, b)
	if err != nil {
		*errs++
	}

	return res
}

func extractRRSet(in []dns.RR, sig *dns.RRSIG) []dns.RR {
	var out []dns.RR

	for _, r := range in {
		if sig.TypeCovered == r.Header().Rrtype {
			if !strings.EqualFold(sig.Header().Name, r.Header().Name) {
				continue
			}

			// Trim TTL
			if r.Header().Ttl > sig.OrigTtl {
				r.Header().Ttl = sig.OrigTtl
			}

			out = append(out, r)
		}
	}
	return out
}
