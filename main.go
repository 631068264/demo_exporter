package main

import (
	"flag"
	"fmt"
	"math"
	"net/http"
	"os"

	"demo_exporter/exporter"
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	dto "github.com/prometheus/client_model/go"
	log "github.com/sirupsen/logrus"
)

const ExporterName = "demo"

func main() {
	var (
		listenAddress = flag.String("web.listen-address", exporter.GetEnv("LISTEN_ADDRESS", ":9121"), "Address to listen on for web interface and telemetry.")
		metricPath    = flag.String("web.telemetry-path", exporter.GetEnv("WEB_TELEMETRY_PATH", "/metrics"), "Path under which to expose metrics.")
	)

	temps := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "pond_temperature_celsius",
		Help:    "The temperature of the frog pond.", // Sorry, we can't measure how badly it smells.
		Buckets: prometheus.LinearBuckets(20, 5, 5),  // 5 buckets, each 5 centigrade wide.
	})

	// Simulate some observations.
	for i := 0; i < 1000; i++ {
		temps.Observe(30 + math.Floor(120*math.Sin(float64(i)*0.1))/10)
	}

	// Just for demonstration, let's check the state of the histogram by
	// (ab)using its Write method (which is usually only used by Prometheus
	// internally).
	metric := &dto.Metric{}
	temps.Write(metric)
	fmt.Println(proto.MarshalTextString(metric))

	log.SetFormatter(&log.TextFormatter{})
	log.SetLevel(log.InfoLevel)
	flag.Parse()
	demoExporter, _ := exporter.NewDemoExporter()
	prometheus.MustRegister(demoExporter)
	http.Handle(*metricPath, promhttp.Handler())
	log.Infof("Providing metrics at %s%s", *listenAddress, *metricPath)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Haproxy Exporter</title></head>
             <body>
             <h1>` + ExporterName + ` Exporter</h1>
             <p><a href='` + *metricPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})
	if err := http.ListenAndServe(*listenAddress, nil); err != nil {
		log.Errorf("Error starting HTTP server %s", err)
		os.Exit(1)
	}
}
