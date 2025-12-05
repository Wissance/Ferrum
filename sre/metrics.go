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
	// Request to dataSource (is not usable for File), cli not produces HttpRequestsCount
	DataSourceRequestsTotalCount *prometheus.CounterVec
	DataSourceRequestDurations   prometheus.Summary
}

func CreateMetricsCollector() *MetricsCollector {
	mc := MetricsCollector{}
	mc.HttpRequestsTotalCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Processed HTTP requests count, partitioned by status",
		},
		[]string{Status})

	mc.HttpRequestsErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_error_total",
			Help: "Processed HTTP requests errors count, partitioned by error type",
		},
		[]string{ErrorType})

	mc.DataSourceRequestsTotalCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "data_source_requests_total",
			Help: "Processed Data source requests count, partitioned by status",
		},
		[]string{Status})

	return &mc
}
