package hnsquery

import (
	"context"
	"encoding/pem"
	"errors"
	"fmt"
	"github.com/imperviousinc/hnsquery/dnssec"
	"github.com/miekg/dns"
	"os"
	"strings"
	"testing"
	"time"
)

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

func TestNewDNSCertVerifier(t *testing.T) {
	r, err := NewResolver(&ResolverConfig{
		Forward: "https://hs.dnssec.dev/dns-query",
	})
	if err != nil {
		t.Fatal(err)
	}

	r.TrustAnchorFunc = func(ctx context.Context, cut string) (*dnssec.Zone, bool, error) {
		if cut == "proofofconcept." {
			fmt.Println("hot lookup")
			time.Sleep(500 * time.Millisecond)
			ds, _ := dns.NewRR("proofofconcept.         21600   IN      DS      60767 15 2 FAF50B8DC0DED5B28E5388F5047805C7417678BE7CAC3AB5DF93823E 9220D87B")
			zone, _ := dnssec.NewZone(cut, []dns.RR{ds})
			zone.Expire.Add(10 * time.Second)
			return zone, true, nil
		}

		if cut == "." {
			return nil, false, errors.New("not supported")
		}

		return nil, false, nil
	}

	verify, err := NewDNSCertVerifier(r)
	if err != nil {
		t.Fatal(err)
	}

	block, _ := pem.Decode([]byte(testCert))
	if block == nil {
		t.Fatal("bad cert")
	}

	ok, err := verify.Verify(context.Background(), &CertVerifyInfo{
		Host:     "proofofconcept",
		Port:     "443",
		Protocol: "tcp",
		RawCerts: [][]byte{block.Bytes},
	})

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("should be secure")
	}

	fmt.Println("done")
}

func TestIntegrationFullVerify(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping")
		return
	}

	client, err := NewClient(&Config{DataDir: os.TempDir()})
	if err != nil {
		t.Fatal(err)
	}
	defer client.Destroy()

	ready := make(chan error)
	client.Start(ready)

	<-ready

	r, err := NewResolver(&ResolverConfig{
		Forward: "https://hs.dnssec.dev/dns-query",
	})
	if err != nil {
		t.Fatal(err)
	}

	r.TrustAnchorFunc = func(ctx context.Context, cut string) (*dnssec.Zone, bool, error) {
		if cut == "proofofconcept." {
			fmt.Println("hot lookup")
			rrs, err := client.GetZone(ctx, strings.TrimSuffix(cut, "."))
			if err != nil {
				return nil, false, err
			}
			var dsSet []dns.RR
			for _, rr := range rrs {
				if rr.Header().Rrtype == dns.TypeDS {
					dsSet = append(dsSet, rr)
				}
			}

			zone, _ := dnssec.NewZone(cut, dsSet)
			zone.Expire.Add(10 * time.Second)
			return zone, true, nil
		}

		if cut == "." {
			return nil, false, errors.New("not supported")
		}

		return nil, false, nil
	}

	verify, err := NewDNSCertVerifier(r)
	if err != nil {
		t.Fatal(err)
	}

	block, _ := pem.Decode([]byte(testCert))
	if block == nil {
		t.Fatal("bad cert")
	}

	now := time.Now()

	ok, err := verify.Verify(context.Background(), &CertVerifyInfo{
		Host:     "proofofconcept",
		Port:     "443",
		Protocol: "tcp",
		RawCerts: [][]byte{block.Bytes},
	})

	fmt.Println("\nElpased: ", time.Since(now).Milliseconds())

	if err != nil {
		t.Fatal(err)
	}

	if !ok {
		t.Fatal("should be secure")
	}

	fmt.Println("done")
}
