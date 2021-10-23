package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	PrometheusRegister = prometheus.NewRegistry()

	ReceivedRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "furby_received_req_total",
			Help: "Total number of received requests.",
		},
	)

	BackendRequests = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "furby_oauth2_token_requests",
			Help: "Total number of backend requests",
		},
	)

	Http500Errors = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "furby_internal_server_errors",
			Help: "Total number of internal server errors",
		},
	)

	RequestTime = prometheus.NewSummary(
		prometheus.SummaryOpts{
			Name: "furby_request_time",
			Help: "Time of request handling",
		},
	)
)

func init() {
	PrometheusRegister.MustRegister(BackendRequests)
	PrometheusRegister.MustRegister(Http500Errors)
	PrometheusRegister.MustRegister(ReceivedRequests)
	PrometheusRegister.MustRegister(RequestTime)
}
