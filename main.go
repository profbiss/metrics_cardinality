package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"

	prom "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
)

func main() {
	http.HandleFunc("/metrics", func(writer http.ResponseWriter, _ *http.Request) {
		url := getEnv("METRICS_URL", "http://localhost:58915/metrics")

		res, err := http.Get(url)
		if err != nil {
			slog.Error("Get metrics failed", "err", err)
			return
		}

		var parser expfmt.TextParser
		mf, err := parser.TextToMetricFamilies(res.Body)
		if err != nil {
			slog.Error("Parse metrics failed", "err", err)
			return
		}

		metricsCardinality := &prom.MetricFamily{
			Name:   Pointer("metrics_cardinality"),
			Help:   Pointer("Metrics label cardinality"),
			Type:   Pointer(prom.MetricType_GAUGE),
			Metric: make([]*prom.Metric, len(mf)),
		}
		var i int
		for _, v := range mf {
			metricsCardinality.Metric[i] = &prom.Metric{
				Label: []*prom.LabelPair{
					{Name: Pointer("metric_name"), Value: v.Name},
				},
				Gauge: &prom.Gauge{
					Value: Pointer(float64(len(v.Metric))),
				},
			}
			i++

			_, err := expfmt.MetricFamilyToText(writer, v)
			if err != nil {
				slog.Error("Write metric failed", "err", err)
				return
			}
		}
		_, err = expfmt.MetricFamilyToText(writer, metricsCardinality)
		if err != nil {
			slog.Error("Write metric failed", "err", err)
			return
		}

		metricsCount := &prom.MetricFamily{
			Name: Pointer("metrics_count"),
			Help: Pointer("Count of metrics"),
			Type: Pointer(prom.MetricType_GAUGE),
			Metric: []*prom.Metric{
				{
					Gauge: &prom.Gauge{
						Value: Pointer(float64(len(mf))),
					},
				},
			},
		}
		_, err = expfmt.MetricFamilyToText(writer, metricsCount)
		if err != nil {
			slog.Error("Write metric failed", "err", err)
			return
		}
	})

	addr := getEnv("LISTEN_ADD", ":8080")
	fmt.Printf("Starting server at addr %v\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		slog.Error("Listen failed", "err", err)
	}
}

// Pointer returns pointer to value.
func Pointer[T any](v T) *T {
	return &v
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
