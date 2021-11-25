package exporter

import (
	"fmt"
	"math/rand"
	"strconv"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

const namespace = "demo"

type DemoExporter struct {
	mutex         sync.RWMutex
	serverMetrics map[string]MetricInfo
	totalScrapes  prometheus.Counter
}

// randFloats [min, max)
func randFloats(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func NewDemoExporter() (*DemoExporter, error) {
	metricCreator := MetricUtil{
		Namespace: namespace,
	}
	return &DemoExporter{
		totalScrapes: metricCreator.NewCounter(
			"exporter_scrapes_total",
			fmt.Sprintf("Current total %s scrapes.", namespace),
		),
		serverMetrics: map[string]MetricInfo{
			"random_value": metricCreator.NewMetric("random_value", "random_value", prometheus.GaugeValue, []string{"demo"}, nil),
		},
	}, nil
}

func (e *DemoExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range e.serverMetrics {
		ch <- m.Desc
	}
	ch <- e.totalScrapes.Desc()
}

func (e *DemoExporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()

	e.totalScrapes.Inc()
	log.Infof("Collect agent")
	for _, m := range e.serverMetrics {
		val := randFloats(1, 100)
		labelValue := strconv.FormatFloat(val, 'E', -1, 64)
		ch <- prometheus.MustNewConstMetric(m.Desc, m.Type, randFloats(1, 100), labelValue)
	}
	ch <- e.totalScrapes
}
