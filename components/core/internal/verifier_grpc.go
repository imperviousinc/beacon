package internal

import (
	"context"
	"errors"

	"github.com/imperviousinc/beacon/components/core/public/proto"
	"github.com/imperviousinc/hnsquery"
	"github.com/miekg/dns"
)

// GRPC Verifier should use a mojo pipe instead.
type CertVerifierGRPC struct {
	proto.UnimplementedCertVerifierServer
	config *Config
}

func (bc *CertVerifierGRPC) VerifyCert(ctx context.Context, req *proto.CertVerifyRequest) (*proto.CertVerifyResponse, error) {
	chain := req.Cert.DerCerts

	// Should never happen
	if len(chain) == 0 {
		return &proto.CertVerifyResponse{
			State: proto.SecurityState_BOGUS,
			Code:  proto.ErrorCode_ERR_TRUST_SERVICE_REQUEST_INVALID,
		}, nil
	}

	leafDer := chain[0]

	// Skip ICANN domains. Consumer of this API should
	// already skip those but just in case.
	h := dns.SplitDomainName(req.Host)
	if len(h) > 0 {
		tld := h[len(h)-1]
		if _, ok := nameConstraints[tld]; ok {
			return &proto.CertVerifyResponse{
				VerifiedCert:   nil,
				State:          proto.SecurityState_INSECURE,
				Code:           proto.ErrorCode_UNKNOWN_ERROR,
				AdditionalInfo: "",
			}, nil
		}
	}

	secure, err := bc.config.verifier.Verify(ctx, &hnsquery.CertVerifyInfo{
		Host:     req.Host,
		Port:     req.Port,
		Protocol: "tcp",
		RawCerts: [][]byte{leafDer},
	})

	// No errors
	if err == nil {
		// Insecure zone
		if !secure {
			return &proto.CertVerifyResponse{
				State: proto.SecurityState_INSECURE,
				Code:  proto.ErrorCode_UNKNOWN_ERROR,
			}, nil
		}
		// DANE verified
		return &proto.CertVerifyResponse{
			State: proto.SecurityState_SECURE,
			Code:  proto.ErrorCode_UNKNOWN_ERROR,
		}, nil
	}

	// Bogus
	switch {
	case errors.Is(err, hnsquery.ErrDNSAuthFailed):
		return &proto.CertVerifyResponse{
			State: proto.SecurityState_BOGUS,
			Code:  proto.ErrorCode_ERR_DNSSEC_PINNED_KEY_NOT_IN_CERT_CHAIN,
		}, nil
	case errors.Is(err, hnsquery.ErrTimeout):
		return &proto.CertVerifyResponse{
			State: proto.SecurityState_BOGUS,
			Code:  proto.ErrorCode_ERR_DNS_TIMED_OUT,
		}, nil
	case errors.Is(err, hnsquery.ErrCancelled):
		return &proto.CertVerifyResponse{
			State: proto.SecurityState_BOGUS,
			// TODO: add more suitable error code
			Code: proto.ErrorCode_ERR_ABORTED,
		}, nil
	case errors.Is(err, hnsquery.ErrNotSynced):
		return &proto.CertVerifyResponse{
			State: proto.SecurityState_BOGUS,
			Code:  proto.ErrorCode_ERR_HNS_IS_SYNCING,
		}, nil
	case errors.Is(err, hnsquery.ErrNoPeers):
		return &proto.CertVerifyResponse{
			State: proto.SecurityState_BOGUS,
			Code:  proto.ErrorCode_ERR_HNS_NO_PEERS,
		}, nil
	}

	return &proto.CertVerifyResponse{
		VerifiedCert:   nil,
		State:          proto.SecurityState_BOGUS,
		Code:           proto.ErrorCode_ERR_DNSSEC_BOGUS,
		AdditionalInfo: "",
	}, nil
}
