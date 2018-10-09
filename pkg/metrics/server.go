package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func healthCheck(res http.ResponseWriter, _ *http.Request) {
	res.WriteHeader(200)
	res.Write([]byte("OK"))
}

func RunServer(port string) error {
	http.HandleFunc("/healthz", healthCheck)
	http.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":"+port, nil)
}
