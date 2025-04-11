// metrics/metrics.go
package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	CertExpiryTime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "certificate_expiry_time_seconds",
			Help: "Expiry time of the certificate in seconds since epoch",
		},
		[]string{"path"},
	)

	CertSubject = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "certificate_subject",
			Help: "Subject of the certificate",
		},
		[]string{"path", "subject"},
	)

	CertIssuer = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "certificate_issuer",
			Help: "Issuer of the certificate",
		},
		[]string{"path", "issuer"},
	)

	CertNotBefore = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "certificate_not_before_seconds",
			Help: "Not Before time of the certificate in seconds since epoch",
		},
		[]string{"path"},
	)

	CertNotAfter = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "certificate_not_after_seconds",
			Help: "Not After time of the certificate in seconds since epoch",
		},
		[]string{"path"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(CertExpiryTime)
	prometheus.MustRegister(CertSubject)
	prometheus.MustRegister(CertIssuer)
	prometheus.MustRegister(CertNotBefore)
	prometheus.MustRegister(CertNotAfter)
}
