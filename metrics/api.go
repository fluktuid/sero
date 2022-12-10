package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	callsReceived = promauto.NewCounter(prometheus.CounterOpts{
		Name: "server_calls_received",
		Help: "The total number of received calls",
	})
	callsFailed = promauto.NewCounter(prometheus.CounterOpts{
		Name: "server_calls_failed",
		Help: "The total number of failed calls",
	})
	callsSuccessful = promauto.NewCounter(prometheus.CounterOpts{
		Name: "server_calls_successful",
		Help: "The total number of processed calls",
	})
)

func RecordRequest() {
	go func() {
		callsReceived.Inc()
	}()
}

func RecordRequestFinish(success bool) {
	go func(success bool) {
		if success {
			callsSuccessful.Inc()
		} else {
			callsFailed.Inc()
		}
	}(success)
}
