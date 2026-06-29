package metrics

import (
	"context"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

var (
	GRPCRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sso",
			Subsystem: "grpc",
			Name:      "requests_total",
			Help:      "total number of grpc requests.",
		},
		[]string{"method", "code"},
	)
	GRPCRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "sso",
			Subsystem: "grpc",
			Name:      "request_duration_seconds",
			Help:      "gRPC request duration in seconds.",
			Buckets:   prometheus.DefBuckets,
		},
		[]string{"method", "code"},
	)

	AuthEventsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: "sso",
			Subsystem: "auth",
			Name:      "events_total",
			Help:      "Total number of auth events.",
		},
		[]string{"event", "status"},
	)
)

func init() {
	prometheus.MustRegister(GRPCRequestsTotal)
	prometheus.MustRegister(GRPCRequestDuration)
	prometheus.MustRegister(AuthEventsTotal)
}

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		start := time.Now()

		resp, err := handler(ctx, req)

		code := status.Code(err).String()
		method := info.FullMethod

		GRPCRequestsTotal.WithLabelValues(method, code).Inc()
		GRPCRequestDuration.WithLabelValues(method, code).Observe(time.Since(start).Seconds())

		return resp, err
	}
}

func IncAuthEvent(event string, status string) {
	AuthEventsTotal.WithLabelValues(event, status).Inc()
}
