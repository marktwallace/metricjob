package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// BulkMessageCounter -- Prometheus Counter
var MainLoopCounter = promauto.NewCounter(prometheus.CounterOpts{
	Name: "main_loop_counter",
	Help: "The number of top level iterations"})
