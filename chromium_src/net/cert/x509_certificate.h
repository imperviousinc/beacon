#ifndef BEACON_CHROMIUM_SRC_NET_CERT_X509_CERTIFICATE_H_
#define BEACON_CHROMIUM_SRC_NET_CERT_X509_CERTIFICATE_H_

#define BEACON_X509_CERT_PROPERTIES bool is_dnssec_cert = false; \
    bool is_hns_hostname = false; 

#include "src/net/cert/x509_certificate.h"

#undef BEACON_X509_CERT_PROPERTIES
#endif // BEACON_CHROMIUM_SRC_NET_CERT_X509_CERTIFICATE_H_
