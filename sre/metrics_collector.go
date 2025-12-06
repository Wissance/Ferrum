package sre

import "github.com/prometheus/client_golang/prometheus"

type StatusLabel string
type ErrorSide string

const (
	Status                  = "status"
	ErrorType               = "error"
	Success     StatusLabel = "Success"
	Fail        StatusLabel = "Fail"
	ClientError ErrorSide   = "ClientError"
	ServerError ErrorSide   = "ServerError"
)

// MetricsCollector is a struct that collects of metrics that is using for application observability
type MetricsCollector struct {
	// Http requests (total both OK and non-OK)
	HttpRequestsTotalCount *prometheus.CounterVec
	// Http requests that were finished with server Error (4xx or 5xx)
	HttpRequestsErrorCount *prometheus.CounterVec
	HttpRequestDurations   prometheus.Summary
	// Request to dataSource (is not usable for File), cli not produces HttpRequestsCount
	DataSourceRequestsTotalCount *prometheus.CounterVec
	DataSourceRequestDurations   prometheus.Summary
}

func CreateMetricsCollector() *MetricsCollector {
	mc := MetricsCollector{}
	// 1. HTTP Section
	mc.HttpRequestsTotalCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Processed HTTP requests count, partitioned by status",
		},
		[]string{Status})
	prometheus.MustRegister(mc.HttpRequestsTotalCount)

	mc.HttpRequestsErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_error_total",
			Help: "Processed HTTP requests errors count, partitioned by error type",
		},
		[]string{ErrorType})
	prometheus.MustRegister(mc.HttpRequestsErrorCount)

	mc.HttpRequestDurations = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "http_request_durations",
			Help:       "Http requests latencies in seconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}})
	prometheus.MustRegister(mc.HttpRequestDurations)
	// 2. DataSource section
	mc.DataSourceRequestsTotalCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "data_source_requests_total",
			Help: "Processed Data source requests count, partitioned by status",
		},
		[]string{Status})
	prometheus.MustRegister(mc.DataSourceRequestsTotalCount)

	mc.DataSourceRequestDurations = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "http_request_durations",
			Help:       "Http requests latencies in seconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}})
	prometheus.MustRegister(mc.DataSourceRequestDurations)

	return &mc
}
