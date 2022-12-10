package metrics

import (
	"net/http"
	"sync/atomic"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func InitAsync() {
	go func() {
		Init()
	}()
}

var isReady = atomic.Bool{}
var isHealthy = atomic.Bool{}

func Init() {
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/readyz", readyZ)
	http.HandleFunc("/healthz", healthZ)
	http.ListenAndServe(":2112", nil)
}

func Ready(ready bool) {
	isReady.Store(ready)
}

func Healthy(healthy bool) {
	isHealthy.Store(healthy)
}

func readyZ(w http.ResponseWriter, _ *http.Request) {
	if isReady.Load() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusTooEarly)
	}
}

func healthZ(w http.ResponseWriter, _ *http.Request) {
	if isHealthy.Load() {
		w.WriteHeader(http.StatusOK)
	} else {
		w.WriteHeader(http.StatusConflict)
	}
}
