package simmet

// SimulatedMetric to hold state
type SimulatedMetric struct {
	Name   string
	Type   string // Counter Gauge Histogram Summary
	Labels []string
	//LabelGenerator          []string
	//ValueGeneratorSelection string
	//ValueGeneratorParams    []float64 //maybe use a set of function
}

// NewSimulatedMetric create
func NewSimulatedMetric() *SimulatedMetric {
	s := &SimulatedMetric{}
	return s
}

// Start app from main
func (s *SimulatedMetric) Start() {
}
