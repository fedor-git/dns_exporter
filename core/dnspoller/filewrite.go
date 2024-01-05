package dnspoller

import (
	"bytes"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

func gatherSpecificMetrics(metricsNames []string, filename string) error {
	registry := prometheus.DefaultGatherer
	metricFamilies, err := registry.Gather()
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	for _, mf := range metricFamilies {
		for _, name := range metricsNames {
			if mf.GetName() == name {
				_, err := expfmt.MetricFamilyToText(&buf, mf)
				if err != nil {
					return err
				}
				break
			}
		}
	}

	return os.WriteFile(filename, buf.Bytes(), 0644)
}