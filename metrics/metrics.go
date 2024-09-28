package metrics

import (
	"maps"
	"slices"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

type MetricCollector interface {
	IncrementCounter(metricName string, opts ...CounterOpts)
}

var Collector *collector

type collector struct{}

type CounterOpts struct {
	Value float64
	Tags  Tags
}

type Tags map[string]string

var countersCache sync.Map

func (r *collector) IncrementCounter(metricName string, opts ...CounterOpts) {
	if len(opts) == 0 {
		opts = []CounterOpts{{Value: 1}}
	}

	labels := prometheus.Labels(opts[0].Tags)
	value, found := countersCache.Load(metricName)
	if !found {
		counterVec := prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: metricName,
			}, slices.Collect(maps.Keys(labels)),
		)
		countersCache.Store(metricName, counterVec)
		if err := prometheus.Register(counterVec); err != nil {
			return
		}
		value = counterVec
	}

	counterVec, ok := value.(*prometheus.CounterVec)
	if !ok {
		return
	}

	metric, err := counterVec.GetMetricWith(labels)
	if err != nil {
		return
	}

	metric.Add(opts[0].Value)
}
