package middleware

import (
	grpcprom "github.com/grpc-ecosystem/go-grpc-middleware/providers/prometheus"
	"github.com/prometheus/client_golang/prometheus"
)

func NewRegistry() *prometheus.Registry {
	return prometheus.NewRegistry()
}

func NewMetrics(registry *prometheus.Registry) *grpcprom.ServerMetrics {
	metrics := grpcprom.NewServerMetrics(
		grpcprom.WithServerHandlingTimeHistogram(
			grpcprom.WithHistogramBuckets([]float64{0.001, 0.01, 0.1, 0.3, 0.6, 1, 3, 6, 9, 20, 30, 60, 90, 120}),
		),
	)
	registry.MustRegister(metrics)
	return metrics
}
