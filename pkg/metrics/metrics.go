package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	ApiVersionLabel = "api_version"
	KindLabel       = "kind"
)

var (
	successRuns = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "shelloperator_success_runs",
		Help: "Number of successful shell executions that have occured, where success is a zero exit code.",
	}, []string{ApiVersionLabel, KindLabel})

	failureRuns = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "shelloperator_failure_runs",
		Help: "Number of failed shell executions that have occured, where failed is a non zero exit code.",
	}, []string{ApiVersionLabel, KindLabel})

	shellRunTime = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name: "shelloperator_run_time_milliseconds",
		Help: "The time it takes for shell executions to run as a histogram in milliseconds.",
	}, []string{ApiVersionLabel, KindLabel})
)

func init() {
	prometheus.MustRegister(successRuns)
	prometheus.MustRegister(failureRuns)
	prometheus.MustRegister(shellRunTime)
}

func IncSuccessRun(apiVersion, kind string) {
	successRuns.With(prometheus.Labels{
		ApiVersionLabel: apiVersion,
		KindLabel:       kind,
	}).Inc()
}

func IncFailureRun(apiVersion, kind string) {
	failureRuns.With(prometheus.Labels{
		ApiVersionLabel: apiVersion,
		KindLabel:       kind,
	}).Inc()
}

func ObserveRunTime(apiVersion, kind string, milliseconds float64) {
	shellRunTime.With(prometheus.Labels{
		ApiVersionLabel: apiVersion,
		KindLabel:       kind,
	}).Observe(milliseconds)
}
