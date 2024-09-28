package main

import (
	"net/http"

	"manc/metrics"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {
	// Expose the registered metrics via HTTP
	http.Handle("/metrics", promhttp.Handler())

	metrics.Collector.IncrementCounter("my_counter", metrics.CounterOpts{Value: 1, Tags: metrics.Tags{"type": "example"}})
	metrics.Collector.IncrementCounter("my_counter", metrics.CounterOpts{Value: 1, Tags: metrics.Tags{"type": "example"}})
	metrics.Collector.IncrementCounter("my_counter", metrics.CounterOpts{Value: 1, Tags: metrics.Tags{"type": "example"}})
	metrics.Collector.IncrementCounter("my_counter", metrics.CounterOpts{Value: 1, Tags: metrics.Tags{"type": "example"}})

	metrics.Collector.IncrementCounter("order_status", metrics.CounterOpts{Value: 1, Tags: metrics.Tags{"status": "pending"}})
	metrics.Collector.IncrementCounter("order_status", metrics.CounterOpts{Value: 1, Tags: metrics.Tags{"status": "completed"}})
	metrics.Collector.IncrementCounter("order_status", metrics.CounterOpts{Value: 1, Tags: metrics.Tags{"status": "cancelled"}})

	metrics.Collector.IncrementCounter("empty")
	metrics.Collector.IncrementCounter("empty")
	metrics.Collector.IncrementCounter("empty")
	metrics.Collector.IncrementCounter("empty")
	metrics.Collector.IncrementCounter("empty")

	if err := http.ListenAndServe(":8081", nil); err != nil {
		panic(err)
	}
}
