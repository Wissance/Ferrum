package sre

import (
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/wissance/Ferrum/utils/httputils"
	"net/http"
	"strings"
)

const (
	pathLabel      = "path"
	statusLabel    = "status"
	errorTypeLabel = "error type"
	clientError    = "client"
	serverError    = "server"
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
		[]string{pathLabel, statusLabel})
	prometheus.MustRegister(mc.HttpRequestsTotalCount)

	mc.HttpRequestsErrorCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_error_total",
			Help: "Processed HTTP requests errors count, partitioned by error type",
		},
		[]string{pathLabel, errorTypeLabel})
	prometheus.MustRegister(mc.HttpRequestsErrorCount)

	mc.HttpRequestDurations = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "http_request_durations",
			Help:       "Http requests latencies in milliseconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}})
	prometheus.MustRegister(mc.HttpRequestDurations)
	// 2. DataSource section
	mc.DataSourceRequestsTotalCount = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "data_source_requests_total",
			Help: "Processed Data source requests count, partitioned by status",
		},
		[]string{statusLabel})
	prometheus.MustRegister(mc.DataSourceRequestsTotalCount)

	mc.DataSourceRequestDurations = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name:       "data_source_request_durations",
			Help:       "Data source requests latencies in milliseconds",
			Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001}})
	prometheus.MustRegister(mc.DataSourceRequestDurations)

	return &mc
}

// HttpMetricsCollectMiddleware function is using to track all HTTP-requests and collect
func (mc *MetricsCollector) HttpMetricsCollectMiddleware(next http.Handler) http.Handler {
	const swaggerPath = "/swagger/"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()
		if strings.Contains(path, swaggerPath) {
			next.ServeHTTP(w, r)
		} else {
			timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
				ms := v * 1000 // make milliseconds
				mc.HttpRequestDurations.Observe(ms)
			}))
			lrw := httputils.NewLoggingResponseWriter(w)
			next.ServeHTTP(lrw, r)
			timer.ObserveDuration()
			// collect request count
			mc.HttpRequestsTotalCount.WithLabelValues(path, http.StatusText(lrw.StatusCode)).Inc()
			if lrw.StatusCode >= 400 {
				if lrw.StatusCode < 500 {
					mc.HttpRequestsErrorCount.WithLabelValues(path, clientError).Inc()
				} else {
					mc.HttpRequestsErrorCount.WithLabelValues(path, serverError).Inc()
				}
			}
		}
	})
}
