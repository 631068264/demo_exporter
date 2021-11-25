package exporter

import "github.com/prometheus/client_golang/prometheus"

type MetricInfo struct {
	Desc *prometheus.Desc
	Type prometheus.ValueType
}

type MetricUtil struct {
	Namespace string
	Subsystem string
}

func (u *MetricUtil) NewMetric(metricName string, docString string, t prometheus.ValueType, labels []string, constLabels prometheus.Labels) MetricInfo {
	return MetricInfo{
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(u.Namespace, u.Subsystem, metricName),
			docString,
			labels,
			constLabels,
		),
		Type: t,
	}
}

func (u *MetricUtil) NewGauge(metricName string, docString string) prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: u.Namespace,
		Name:      metricName,
		Help:      docString,
	})
}
func (u *MetricUtil) NewCounter(metricName string, docString string) prometheus.Counter {
	return prometheus.NewCounter(prometheus.CounterOpts{
		Namespace: u.Namespace,
		Name:      metricName,
		Help:      docString,
	})
}
