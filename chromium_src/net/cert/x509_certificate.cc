#define BEACON_X509_CERT_PICKLE_READ auto cert = CreateFromDERCertChainUnsafeOptions(cert_chain, options);  \
    if (!pickle_iter->ReadBool(&cert->is_dnssec_cert)) return nullptr;                               \
    if (!pickle_iter->ReadBool(&cert->is_hns_hostname)) return nullptr;                               \
    return cert;

#define BEACON_X509_CERT_PICKLE_PERSIST pickle->WriteBool(is_dnssec_cert); \
    pickle->WriteBool(is_hns_hostname);

#include "src/net/cert/x509_certificate.cc"

#undef BEACON_X509_CERT_PICKLE_READ
#undef BEACON_X509_CERT_PICKLE_PERSIST
