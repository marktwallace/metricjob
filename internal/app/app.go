package app

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"time"

	"github.com/caarlos0/env/v6"
	"github.com/marktwallace/metricjob/internal/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// AppContext to hold environment and other resources
type AppContext struct {
	RunTimeInSeconds int64 `env:"RUN_TIME_IN_SECONDS,required"`
	// The following four strings are comma seperated lists
	CounterName  []string  `env:"COUNTER_NAME"`
	CounterRatio []float64 `env:"COUNTER_RATIO"`

	CounterVecName        []string  `env:"COUNTER_VEC_NAME"`
	CounterVecLabel       []string  `env:"COUNTER_VEC_LABEL"`
	CounterVecCardinality []int     `env:"COUNTER_VEC_CARDINALITY"`
	CounterVecRatio       []float64 `env:"COUNTER_VEC_RATIO"`

	GaugeVecName              []string  `env:"GAUGE_VEC_NAME"`
	GaugeVecLabel             []string  `env:"GAUGE_VEC_LABEL"`
	GaugeVecCardinality       []int     `env:"GAUGE_VEC_CARDINALITY"`
	GaugeVecMin               []float64 `env:"GAUGE_VEC_MIN"`
	GaugeVecMax               []float64 `env:"GAUGE_VEC_MAX"`
	GaugeVecNoise             []float64 `env:"GAUGE_VEC_NOISE"`
	GaugeVecPeriodSeconds     []int     `env:"GAUGE_VEC_PERIOD_SECONDS"`
	GaugeVecPhaseShiftSeconds []int     `env:"GAUGE_VEC_PHASE_SHIFT_SECONDS"`

	// internal state
	MapToPrometheusCounter    map[string]prometheus.Counter
	MapToPrometheusCounterVec map[string]*prometheus.CounterVec
	MapToPrometheusGaugeVec   map[string]*prometheus.GaugeVec
}

// NewApp create context
func NewApp() *AppContext {
	app := &AppContext{}
	if err := env.Parse(app); err != nil {
		log.Panicln(err)
	}

	app.MapToPrometheusCounter = make(map[string]prometheus.Counter)
	app.MapToPrometheusCounterVec = make(map[string]*prometheus.CounterVec)
	app.MapToPrometheusGaugeVec = make(map[string]*prometheus.GaugeVec)

	return app
}

func (app *AppContext) GetPrometheusCounter(name string) prometheus.Counter {
	_, exists := app.MapToPrometheusCounter[name]
	if !exists {
		app.MapToPrometheusCounter[name] =
			promauto.NewCounter(prometheus.CounterOpts{Name: name})
	}
	return app.MapToPrometheusCounter[name]
}

func (app *AppContext) GetPrometheusCounterVec(name string, label string) *prometheus.CounterVec {
	key := name + label
	_, exists := app.MapToPrometheusCounterVec[key]
	if !exists {
		app.MapToPrometheusCounterVec[key] =
			promauto.NewCounterVec(prometheus.CounterOpts{Name: name}, []string{label})
	}
	return app.MapToPrometheusCounterVec[key]
}

func (app *AppContext) GetPrometheusGaugeVec(name string, label string) *prometheus.GaugeVec {
	key := name + label
	_, exists := app.MapToPrometheusGaugeVec[key]
	if !exists {
		app.MapToPrometheusGaugeVec[key] =
			promauto.NewGaugeVec(prometheus.GaugeOpts{Name: name}, []string{label})
	}
	return app.MapToPrometheusGaugeVec[key]
}

func (app *AppContext) updateCounters() {
	if len(app.CounterName) == 0 {
		return
	}
	for i, name := range app.CounterName {
		if app.CounterRatio[i] > rand.Float64() {
			counter := app.GetPrometheusCounter(name)
			counter.Inc()
		}
	}
}

func (app *AppContext) updateCounterVecs() {
	if len(app.CounterVecName) == 0 {
		return
	}
	for i, name := range app.CounterVecName {
		label := app.CounterVecLabel[i]
		vec := app.GetPrometheusCounterVec(name, label)
		for k := 0; k < app.CounterVecCardinality[i]; k++ {
			if app.CounterVecRatio[i] > rand.Float64() {
				vec.With(prometheus.Labels{label: strconv.Itoa(k)}).Add(1)
			}
		}
	}
}

func (app *AppContext) updateGaugeVecs(t time.Time) {
	if len(app.GaugeVecName) == 0 {
		return
	}
	for i, name := range app.GaugeVecName {
		label := app.GaugeVecLabel[i]
		vec := app.GetPrometheusGaugeVec(name, label)
		for k := 0; k < app.GaugeVecCardinality[i]; k++ {
			period := app.GaugeVecPeriodSeconds[i]
			shift := k * app.GaugeVecPhaseShiftSeconds[i]
			phaseSec := (t.Unix() + int64(shift)) % int64(period)
			phase := (2 * math.Pi * float64(phaseSec)) / float64(period)
			amplitude := (app.GaugeVecMax[i] - app.GaugeVecMin[i]) / 2.0
			mean := (app.GaugeVecMax[i] + app.GaugeVecMin[i]) / 2.0
			val := math.Sin(phase)*amplitude + mean
			val += rand.Float64() * app.GaugeVecNoise[i]
			vec.With(prometheus.Labels{label: strconv.Itoa(k)}).Set(val)
			//fmt.Printf("%d %6.4f %6.2f %s \n", k, phase, val, strings.Repeat("*", int(val)))
		}
	}
}

func (app *AppContext) updateMetrics(t time.Time) {
	app.updateCounters()
	app.updateCounterVecs()
	app.updateGaugeVecs(t)
}

func (app *AppContext) tickerLoop(ticker *time.Ticker, done chan bool) {
	for {
		select {
		case <-done:
			fmt.Println("Done")
			return
		case t := <-ticker.C:
			metrics.MainLoopCounter.Inc()
			app.updateMetrics(t)
			fmt.Println("Tick at", t)
		}
	}
}

// Start app from main
func (app *AppContext) Start() {
	metrics.MainLoopCounter.Inc()

	ticker := time.NewTicker(10 * time.Second)
	done := make(chan bool)
	go app.tickerLoop(ticker, done)

	time.Sleep(time.Duration(app.RunTimeInSeconds) * time.Second)
	ticker.Stop()
	done <- true
	fmt.Println("Ticker stopped")
}
