package dnssec

import (
	"bufio"
	"context"
	"errors"
	"github.com/miekg/dns"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"testing"
	"time"
)

type testHDR struct {
	verifyDNSKeys bool
	zone          string
	time          time.Time
	dnsKeys       string
	anchors       string
	minRSA        int
}

type testCase struct {
	inputMsg    string
	filteredMsg string
	name        string
	time        time.Time
	secure      bool
	bogus       bool
}

func TestVerify(t *testing.T) {
	wd, _ := os.Getwd()

	dir := path.Join(wd, "testdata")
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, file := range files {
		if !strings.HasPrefix(file.Name(), "val_") {
			continue
		}

		f, err := os.Open(path.Join(dir, file.Name()))
		if err != nil {
			t.Fatal(err)
		}

		parseTestFile(t, f, runTest)
	}
}

func Test_filterDigest(t *testing.T) {
	dsSet := `
ns.forever. 0 IN DS 21761 13 1 004ad4545ffc78c2a19853b4dd5b6b1db96e1c8a
ns.forever. 0 IN DS 21761 13 2 3b606b0aff27ad10e5e8903d5bf3cd36f7abca44d1ba6f9c59099372936a9845
ns.forever. 0 IN DS 21761 13 4 1f20aed75154bd655ac7f9693744afa8a2e97dcffaec0bfd22363d5f70969292cadd001630e9962404a51bb51e35d07e
`
	want := `
ns.forever. 0 IN DS 21761 13 4 1f20aed75154bd655ac7f9693744afa8a2e97dcffaec0bfd22363d5f70969292cadd001630e9962404a51bb51e35d07e
`

	rrs := zoneToRecords(dsSet)
	set, err := filterDS("ns.forever.", rrs)
	if err != nil {
		t.Fatal(err)
	}

	got := recordsToZone(dsToRR(set))
	want = recordsToZone(zoneToRecords(want))

	if got != want {
		t.Fatalf("got dsSet = %s \n want = %s", got, want)
	}

}

func Test_canonicalNameCompare2(t *testing.T) {
	// same tests used by unbound
	tests := []struct {
		name1  string
		name2  string
		result int
	}{
		{
			"",
			"",
			0,
		},
		{
			"example.net.",
			"example.net.",
			0,
		},
		{
			"test.example.net.",
			"test.example.net.",
			0,
		},
		{
			"com.",
			"",
			1,
		},
		{
			"",
			"com.",
			-1,
		},
		{
			"example.com.",
			"com.",
			1,
		},
		{
			"com.",
			"example.com.",
			-1,
		},
		{
			"example.com.",
			"",
			1,
		},
		{
			"",
			"example.com.",
			-1,
		},
		{
			"com.",
			"net.",
			-1,
		},
		{
			"net.",
			"com.",
			1,
		},
		{
			"net.",
			"org.",
			-1,
		},
		{
			"neta.",
			"net.",
			1,
		},
		{
			"ne.",
			"neta.",
			-1,
		},
		{
			"aag.example.net.",
			"bla.example.net.",
			-1,
		},
	}

	for _, test := range tests {
		t.Run(test.name1+","+test.name2, func(t *testing.T) {
			v, err := canonicalNameCompare(test.name1, test.name2)
			if err != nil {
				t.Fatal(err)
			}
			if v != test.result {
				t.Fatalf("got result = %v want %v", v, test.result)
			}
		})
	}
}

func Test_canonicalNameCompare(t *testing.T) {
	zone := `\001.z.example. 300 IN A 127.0.0.1
\200.z.example. 300 IN A 127.0.0.1
example.    300 IN A 127.0.0.1
a.example.  300 IN A 127.0.0.1
yljkjljk.a.example. 300 IN A 127.0.0.1
Z.a.example. 300 IN A 127.0.0.1
zABC.a.EXAMPLE. 300 IN A 127.0.0.1
*.z.example. 300 IN A 127.0.0.1
z.example. 300 IN A 127.0.0.1
`
	want := `example.	300	IN	A	127.0.0.1
a.example.	300	IN	A	127.0.0.1
yljkjljk.a.example.	300	IN	A	127.0.0.1
Z.a.example.	300	IN	A	127.0.0.1
zABC.a.EXAMPLE.	300	IN	A	127.0.0.1
z.example.	300	IN	A	127.0.0.1
\001.z.example.	300	IN	A	127.0.0.1
*.z.example.	300	IN	A	127.0.0.1
\200.z.example.	300	IN	A	127.0.0.1
`

	rrs := zoneToRecords(zone)
	sort.Slice(rrs, func(i, j int) bool {
		res, err := canonicalNameCompare(rrs[i].Header().Name, rrs[j].Header().Name)
		if err != nil {
			panic(err)
		}

		return res < 0
	})

	got := recordsToZone(rrs)
	if got != want {
		t.Fatalf("got zone = \n`%s`\nwant zone = \n`%s`\n", got, want)
	}

	if res, err := canonicalNameCompare("", "."); err != nil || res != 0 {
		t.Fatal("failed to compare root label")
	}

	if res, _ := canonicalNameCompare("\\001.", "."); res != 1 {
		t.Fatal("root label is smaller than other labels")
	}

	if res, err := canonicalNameCompare("eXampl\\069.", "exam\\112le"); err != nil || res != 0 {
		t.Fatal("comparison must be case/formatting insensitive")
	}

	if res, _ := canonicalNameCompare(strings.Repeat("a", 63)+".example", "example"); res != 1 {
		t.Fatal("could not compare long labels")
	}

	if _, err := canonicalNameCompare(strings.Repeat("a", 64)+".example", "example"); err == nil {
		t.Fatal("got no error on invalid label length")
	}
}

func sectionsMatch(t *testing.T, name string, a, b []dns.RR) {
	secA := recordsToZoneSorted(a)
	secB := recordsToZoneSorted(b)

	if secA != secB {
		t.Fatalf("%s section dont't match \n%s != \n%s", name, secA, secB)
	}
}

func runTest(t *testing.T, hdr *testHDR, tc *testCase) {
	dsSet := zoneToRecords(hdr.anchors)
	dnskeyMsg, err := stringToMsg(hdr.dnsKeys)
	if err != nil {
		t.Fatal(err)
	}

	var keys map[uint16]*dns.DNSKEY
	verifyMessage := false

	if strings.TrimSpace(tc.filteredMsg) != "" {
		verifyMessage = true
	}

	filtered, err := stringToMsg(tc.filteredMsg)
	if err != nil {
		t.Fatal(err)
	}

	if hdr.verifyDNSKeys {
		t.Run("verify dnskeys", func(t *testing.T) {
			var err error

			z, err := NewZone(hdr.zone, dsSet)
			if err != nil {
				t.Fatal(err)
			}

			z.MinRSA = hdr.minRSA
			z.CurrentTime = hdr.time

			keys, err = z.VerifyDNSKeys(dnskeyMsg)
			if err != nil {
				t.Fatal(err)
			}
		})

	} else {
		keys = make(map[uint16]*dns.DNSKEY)

		for _, rr := range dnskeyMsg.Ns {
			key := rr.(*dns.DNSKEY)
			keys[key.KeyTag()] = key
		}
	}

	testMsg, err := stringToMsg(tc.inputMsg)
	if err != nil {
		t.Fatal(err)
	}

	currTime := hdr.time
	if !tc.time.IsZero() {
		currTime = tc.time
	}

	z := Zone{
		Name:        hdr.zone,
		Keys:        keys,
		CurrentTime: currTime,
		MinRSA:      hdr.minRSA,
	}
	ok, err := z.Verify(context.Background(), testMsg, testMsg.Question[0].Name, testMsg.Question[0].Qtype)
	if tc.bogus {
		if err == nil {
			t.Fatalf("got no error, want bogus")
		}
	} else if err != nil {
		t.Fatal(err)
	} else if tc.secure != ok {
		t.Fatalf("got secure = %v, want %v", ok, tc.secure)
	}

	if verifyMessage {
		sectionsMatch(t, "authority", testMsg.Ns, filtered.Ns)
		sectionsMatch(t, "answer", testMsg.Answer, filtered.Answer)
		sectionsMatch(t, "additional", testMsg.Extra, filtered.Extra)
	}
}

func parseTestFile(t *testing.T, f *os.File, run func(t *testing.T, hdr *testHDR, tc *testCase)) {
	sc := bufio.NewScanner(f)

	var begin bool

	var th testHDR
	var tc testCase

	scanning := -1

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		switch {
		case strings.HasPrefix(line, "[ZONE]"):
			th.zone = ""
			parseKeyValPairs(line[6:], ",", func(key string, val string) {
				if key == "origin" {
					th.zone = val
				}

				if key == "time" {
					t, err := dns.StringToTime(val)
					if err != nil {
						panic("unable to parse time used in test: " + err.Error())
					}

					th.time = time.Unix(int64(t), 0)
				}
			})
			continue
		case line == "[TRUST_ANCHORS]":
			th.anchors = ""
			scanning = 1
			continue
		case strings.HasPrefix(line, "[DNSKEYS]"):
			th.dnsKeys = ""
			scanning = 2
			th.minRSA = DefaultMinRSAKeySize
			th.verifyDNSKeys = true

			parseKeyValPairs(line[9:], ",", func(key string, val string) {
				if key == "min_rsa_keysize" {
					var err error
					if th.minRSA, err = strconv.Atoi(val); err != nil {
						t.Fatal(err)
					}
				}
				if key == "verify" {
					th.verifyDNSKeys = val == "1"
				}
			})

			continue
		case line == "[INPUT]":
			scanning = 3
			continue
		case line == "[VERIFY_MESSAGE]":
			scanning = 4
			continue
		case strings.HasPrefix(line, "[TEST_BEGIN]"):
			if begin {
				panic("test didn't end")
			}
			begin = true
			parseKeyValPairs(line[12:], ",", func(key string, val string) {
				if key == "name" {
					tc.name = val
				}

				if key == "time" {
					t, err := dns.StringToTime(val)
					if err != nil {
						panic("unable to parse time used in test: " + err.Error())
					}

					tc.time = time.Unix(int64(t), 0)
				}
			})
			continue
		case strings.HasPrefix(line, "[TEST_END]"):
			begin = false
			t.Run(th.zone+":"+tc.name, func(t *testing.T) {
				run(t, &th, &tc)
			})
			tc = testCase{}
			continue
		case strings.HasPrefix(line, "[RESULT]"):
			parseKeyValPairs(line[8:], ",", func(key string, val string) {
				switch key {
				case "secure":
					tc.secure = val == "1"
				case "bogus":
					tc.bogus = val == "1"
				}
			})
			continue
		}

		switch scanning {
		case 1:
			th.anchors += line + "\n"
		case 2:
			th.dnsKeys += line + "\n"
		case 3:
			tc.inputMsg += line + "\n"
		case 4:
			tc.filteredMsg += line + "\n"
		}
	}
}

func zoneToRecords(z string) []dns.RR {
	var records []dns.RR
	tokens := dns.NewZoneParser(strings.NewReader(z), "", "")
	for x, ok := tokens.Next(); ok; x, ok = tokens.Next() {
		err := tokens.Err()
		if err != nil {
			panic(err)
		}
		records = append(records, x)
	}
	return records
}

func recordsToZone(rrs []dns.RR) string {
	var b strings.Builder
	for _, rr := range rrs {
		b.WriteString(rr.String() + "\n")
	}
	return b.String()
}

func recordsToZoneSorted(rrs []dns.RR) string {
	sort.Slice(rrs, func(i, j int) bool {
		return rrs[i].String() < rrs[j].String()
	})

	return recordsToZone(rrs)
}

func stringToMsg(str string) (*dns.Msg, error) {
	msg := new(dns.Msg)
	sc := bufio.NewScanner(strings.NewReader(str))
	sectionId := -2

	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if strings.HasPrefix(line, ";") {
			line := strings.TrimFunc(line, func(r rune) bool {
				if r == ';' || r == ' ' || r == '\t' {
					return true
				}
				return false
			})

			switch {
			case strings.HasPrefix(line, "->>HEADER<<-"):
				line := strings.TrimSpace(line[12:])
				err := parseKeyValPairs(line, ",", func(key string, val string) {
					switch key {
					case "opcode":
						msg.Opcode = dns.StringToOpcode[val]
					case "status":
						msg.Rcode = dns.StringToRcode[val]
					}
				})
				if err != nil {
					return nil, err
				}
			case strings.HasPrefix(line, "flags"):
				err := parseKeyValPairs(line, ";", func(key string, val string) {
					switch key {
					case "flags":
						flags := strings.Split(val, " ")
						for _, flag := range flags {
							flag = strings.TrimSpace(flag)
							switch flag {
							case "qr":
								msg.Response = true
							case "rd":
								msg.RecursionDesired = true
							case "ra":
								msg.RecursionAvailable = true
							case "ad":
								msg.AuthenticatedData = true
							case "cd":
								msg.CheckingDisabled = true
							case "aa":
								msg.AuthenticatedData = true
							case "tc":
								msg.Truncated = true
							}
						}
					}
				})
				if err != nil {
					return nil, err
				}
			case strings.Contains(line, "flags: do"):
				msg.SetEdns0(4096, true)
			case line == "QUESTION SECTION:":
				sectionId = -1
			case line == "ANSWER SECTION:":
				sectionId = 0
			case line == "AUTHORITY SECTION:":
				sectionId = 1
			case line == "ADDITIONAL SECTION:":
				sectionId = 2
			case sectionId == -1:
				parts := strings.Fields(line)
				if len(parts) != 3 {
					return nil, errors.New("bad question")
				}

				msg.Question = make([]dns.Question, 1)
				msg.Question[0] = dns.Question{
					Name:   parts[0],
					Qtype:  dns.StringToType[parts[2]],
					Qclass: dns.StringToClass[parts[1]],
				}
			}
			continue
		}

		if line == "" {
			continue
		}

		switch sectionId {
		case 0:
			rr, err := dns.NewRR(line)
			if err != nil {
				return nil, err
			}
			msg.Answer = append(msg.Answer, rr)
		case 1:
			rr, err := dns.NewRR(line)
			if err != nil {
				return nil, err
			}
			msg.Ns = append(msg.Ns, rr)
		case 2:
			rr, err := dns.NewRR(line)
			if err != nil {
				return nil, err
			}
			msg.Extra = append(msg.Extra, rr)
		}
	}

	return msg, nil
}

func parseKeyValPairs(line string, sep string, onRead func(key string, val string)) error {
	parts := strings.Split(line, sep)
	for _, part := range parts {
		keyVal := strings.Split(part, ":")
		if len(keyVal) < 2 {
			return errors.New("bad key value pair")
		}
		key := strings.TrimSpace(keyVal[0])
		val := strings.TrimSpace(keyVal[1])

		if len(keyVal) > 2 {
			key = ""
			val = strings.Join(keyVal, ":")
		}

		onRead(key, val)
	}
	return nil
}

func dsToRR(dsSet []*dns.DS) []dns.RR {
	var rrs []dns.RR
	for _, rr := range dsSet {
		rrs = append(rrs, rr)
	}
	return rrs
}
