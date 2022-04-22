package mobile

import (
	"encoding/pem"
	"fmt"
	"math/rand"
	"strconv"
	"sync"
	"testing"
	"time"
)

// cert name: proofofconcept
const testCert = `-----BEGIN CERTIFICATE-----
MIIFNjCCAx6gAwIBAgIUNGjEZ7cD77NfEr71EbJU3qR3x4AwDQYJKoZIhvcNAQEL
BQAwGzEZMBcGA1UEAwwQKi5wcm9vZm9mY29uY2VwdDAeFw0yMTExMDUwMDAwMzha
Fw0yMjExMDUwMDAwMzhaMBsxGTAXBgNVBAMMECoucHJvb2ZvZmNvbmNlcHQwggIi
MA0GCSqGSIb3DQEBAQUAA4ICDwAwggIKAoICAQDINoydwO3KeVJqK1b1vvPC+sRR
Iu0qTKv7PxlRIJ3d+FCaC/ZOT3lmZWhXUuPIY7EsZM+GizyQW5eeO6hkdlAQVCdG
Fayvd/Xas5PIcpZSyQdcnSfvr8lnM8Xk3aVRXlGT1Np2cwqkYcURYrODxDIKdQox
g4qBADFzZNi1y11w4XZabmK1M03COS7yScIKKb8mnhtvB4iP77n+cIBwPfs2DdDf
9Mt6slR/y4mR+w54mn9EFLQTRm9ngwLsevBtqt4iffpvhBVspBSIIHV8yA0NvjIy
pslZZfpC7f9Qndn0bealZVjLZcx1eNkSZWRcGMNW2TyhYE4+Zi14ETHjPIhyEsIA
0HDwWc+SbpJOroRoy9SjCNmuq9fe3YMZRoz+/gcyTmpAgaVa9FU2ATkFhZFPpNmf
GF5WSXaijSLfvZCKCzkGb9DahDBLDy/Xyycwqx1Y8KHya9E56xSK4nlPyOsd098E
5YZ31ChOVjCeSAHeFy7SAfIXTeYb5CJ3vzBmYXCe6hESBjy16E3DFMQmKm3viS+6
rQzaBheAHOzxiyEvh3iKMX3utXJySW+01IsL6tdRtwzZQJcCS/MjOBFyAc8aV5BF
NjQjq1qjlgMpFCM5dBzryjMJs5J4Rczra1Vd7woXiAmfz+syFIObdmAHmXD/AlHW
lUK8a4ncAZLup7+YLQIDAQABo3IwcDAOBgNVHQ8BAf8EBAMCBaAwEwYDVR0lBAww
CgYIKwYBBQUHAwEwDAYDVR0TAQH/BAIwADA7BgNVHREENDAygg5wcm9vZm9mY29u
Y2VwdIIQKi5wcm9vZm9mY29uY2VwdIIOMTY1LjIyNy45My4xMTcwDQYJKoZIhvcN
AQELBQADggIBAH+MvnWmz6yuAqN4CskJTnUp0DeI94XQRsQc9km6vMnyh3gKRwEB
QILyG/RD4eOFrQQ4mhWxSm9Q9VsIOIVaGfbsK5QQWXpVnPJ7mDEQUrbCA/E0iFUV
+8p4wI4d5GmC2fpdVMlT6tTCGrg2B4UaENd+5NrYi+WISWzjIS5k73dLzIN5V2qG
3wJyetonGKnc4iT0FGpG2oFgwtXgSbBjxmFVaflJ29c8MJqoAzK/9ZH0DxPThGGl
EeqvCNJcRB1ekqauS+gI/n/v1liVFwKZSengZ7jYAUGiTFwzmt9OlSsrD3QvQBuO
yHThswrWTlgeiHSlOYKoECSufn5Ru8BjaTCCW5B5UGxDF8bvPqU4j/Svx6jbEHx4
Gizr6rlgatRJXzcWRKDH0YnELCAQOWan5J/C+5pRuF00mm/cDA3DKZs0dYKGYGIA
Ak+9aCDfmUgseRmRf51eYui4TggVBoXZYA8j0Hz6wPp2zdUqexfS7xuJwEKwt3Ev
gnXlXj3WkAzYfaTMFK5DC8Ys3DMPZk+E5m1YRN7wXhvwpw+MKvXx+eCzbcINmnGY
/SoKPk5ffWnz2uns9FtvIpr1QltBoIA96QfQkvbkm7S2H2EjKGq7ajrDbpGmkNFt
qdycsDX7VrZb6oug6X9siewjN5FpNEI+wqhr4NFBLr4kclt6nhj7/pha
-----END CERTIFICATE-----`

// cert name: welcome.nb
const testCert2 = `-----BEGIN CERTIFICATE-----
MIIFFjCCAv6gAwIBAgIUGUIGnbuVrb9Dv2Ablgth0XRD55MwDQYJKoZIhvcNAQEL
BQAwFzEVMBMGA1UEAwwMKi53ZWxjb21lLm5iMB4XDTIxMTExODEwNTc1MloXDTIy
MTExODEwNTc1MlowFzEVMBMGA1UEAwwMKi53ZWxjb21lLm5iMIICIjANBgkqhkiG
9w0BAQEFAAOCAg8AMIICCgKCAgEAtQJEflxC9q2ut/Bf5aakgQu2XVEmJu7uIEiZ
Ui25xp1m6zZS+OLo0gWp5KOzxkltLzD1b0ZuqWeVBg8F31OBm1RFCB5DtprxhaTU
zTM1iP8BQaHNoyOXA/cAudz7y7oLACWGhgkP+jl8cNKkooA8jRjFkhgZedqx7ipZ
CfcbmfjjtGXUnzdG6ymKakuvJhnsKWCiiYUYvSlPc5Y7zl1zwwV7Wy8CRRifN6U1
zYDxmlyKp4Y2Oe+9cK/1j23G7gJmd8EPRNtoc0pYwc7Eiz54pf2KU8pOUXitkIDt
S6tC9YTpn8mK7cn/uu6L/E2Mx7ACvBLrK7JsR36LVjztbxYt6UeQJhhpa0uuXZNN
hD363shksCp3A5Tt4nfB9rxX5u3NC1+w5SoZkqMZGiHzdyYN7pAEXhDqpm+eJJo1
uqjNUkNA3VFHZrZA1/kPK5XkLbTHY+MoyhHAG6Kot6zLPlUyrYM25lFjOHy12154
IxrHAQez8PHV9EsQgef9iGI8WBmcfogJqGXUWccdk3I+NL9JmybJm9epWwrvgPRU
M3tjCSEREoTTLKxQ9RoZfDdOuCq573WNOenhRErQmU+jb6YcSgyEV2iUnYMoJ3Uf
A3CwN86h/AyRBMpc/gy75nXID/TBNk7p0cRK1tErJ4PVbAZg3sRX3gBA3wSOK0VS
ufqRLV0CAwEAAaNaMFgwDgYDVR0PAQH/BAQDAgWgMBMGA1UdJQQMMAoGCCsGAQUF
BwMBMAwGA1UdEwEB/wQCMAAwIwYDVR0RBBwwGoIKd2VsY29tZS5uYoIMKi53ZWxj
b21lLm5iMA0GCSqGSIb3DQEBCwUAA4ICAQAQ7TS7cyzimV1Nca0M6qKoH3/5Ty+H
2DWJxYaaBg0PButrNwFirUml5QLpYwE0beg8i8P8pAo5c2BVV+eR/ALQBMnCKdLE
FILabL9RaXsYx/yNEA6auOo/zvASj+I+piSKRaWJJoNtPaEgEpYx7JZtkoq30w+y
Di7kZ6huLfagSwJ+4t3AspbgV/x7qq/D8NpAsWQ+wF9nuTBw0VId/050wbTg0w/r
15+48SF2KolIxo+vrl+iydfuvfRIeaQqCNIGTmdrdVZl8RMr79AOQYmh0vnlk+b0
jQF5yLNvHQqWNrwCmBOK4CSZF3eTcHBHEeVti/PVg4s8uF0HtLRAzAN4JZRTpGU8
zAx2n2fAXEh3OZn8/GNH7fS0HJMZkAmLRB6txZmtIzGhMG2KwACf1hqf15dj0XNz
mBgPyWldQPqOYF+Lwi+ZfhGxT+WNwYlHB7ZRu9sDdXOno9pRroJi4ReSx60epOCV
inh4Sz+ohtPh/dxOHHVpqBUwhjY41AnKwA3Nf8KBLg96lWMpONTwcw927nnoK8ih
BRGZ5CRy96LkFpg2tFavuM4/GQ7T1h2DQTCgRORBRf2Dr+WVwmAcgmeYdA/6q+N0
o0jQRW/R09O/96UV8R1dFKU6Cd3xLcRBgKBEjlAr/rqEwdnbXlh0hRj3iIgpf4QF
Dkvzk4cDrPbJTQ==
-----END CERTIFICATE-----`

// TODO: write proper tests
func TestNewVerifier(t *testing.T) {
	if testing.Short() {
		t.Fatal("skipping")
		return
	}

	v, err := NewVerifier("https://hns.dnssec.dev/dns-query")
	if err != nil {
		t.Fatal(err)
	}

	go v.LaunchTA()
	defer v.ShutdownTA()

	block, _ := pem.Decode([]byte(testCert))
	if block == nil {
		t.Fatal("bad cert")
	}

	v.tldDiskCache.cleanUp(true)

	result := v.VerifyCert(block.Bytes, "443", "tcp", "proofofconcept")
	if result != HNSNotSynced && result != HNSNoPeers {
		t.Error("client should still be syncing or finding peers")
		return
	}

	ticker := time.NewTicker(time.Second)

Wait:
	for {
		select {
		case <-ticker.C:
			if v.Ready() && v.ActivePeerCount() > 0 {
				break Wait
			}
		}
	}

	result = v.VerifyCert(block.Bytes, "443", "tcp", "proofofconcept")
	if result != HNSSecure {
		t.Error("want secure")
		return
	}

	result = v.VerifyCert(block.Bytes, "443", "tcp", "letsdane")
	if result != HNSBogus {
		t.Error("want bogus")
		return
	}

	// should fail name checks
	result = v.VerifyCert(block.Bytes, "443", "tcp", "welcome.nb")
	if result != HNSBogus {
		t.Error("want bogus bad name")
		return
	}

	// valid cert name
	block, _ = pem.Decode([]byte(testCert2))
	if block == nil {
		t.Fatal("bad cert")
	}

	result = v.VerifyCert(block.Bytes, "443", "tcp", "welcome.nb")
	if result != HNSInsecure {
		t.Error("want insecure")
		return
	}

	// disable name checks
	v.disableNameChecks = true

	result = v.VerifyCert(block.Bytes, "443", "tcp", "hns.blockclock")
	if result != HNSBogus {
		t.Error("want bogus bad cert")
		return
	}

	names := []string{"3b",
		"proofofconcept",
		"example.com",
		"schematic",
		"howtomakepancakes",
		"nb",
		"humbly",
		"niami",
		"hns.blockclock",
		"forever"}

	var wg sync.WaitGroup
	wg.Add(len(names) + 30)

	// Some concurrency
	for i := 0; i < 30; i++ {
		name := "t" + strconv.Itoa(rand.Int()) + "test" + strconv.Itoa(rand.Int()) + "."
		go func() {
			defer wg.Done()
			result := v.VerifyCert(block.Bytes, "443", "tcp", name)
			if result == HNSSecure {
				t.Error("want result != secure")
				return
			}
			fmt.Printf("Verify %s: result %d\n", name, result)
		}()
	}
	time.Sleep(2 * time.Second)
	for _, name := range names {
		go func(name string) {
			defer wg.Done()
			result := v.VerifyCert(block.Bytes, "885", "tcp", name)
			if result == HNSSecure {
				t.Error("want result != secure")
				return
			}
			fmt.Printf("Verify %s: result %d\n", name, result)
		}(name)
	}

	wg.Wait()
}
