package dnspoller

import (
	"github.com/fedor-git/dns_exporter/core/config"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type MetricDescription struct {
    Name        string
    Help        string
    LabelNames  []string
}

var (
    dnsLookupTimeMetric = MetricDescription{
        Name:       "dns_lookup_time",
        Help:       "DNS lookup time in seconds",
        LabelNames: []string{"dns_server", "hostname"},
    }

    dnsAvailabilityMetric = MetricDescription{
        Name:       "dns_availability",
        Help:       "DNS server availability",
        LabelNames: []string{"dns_server"},
    }
)

type DNSPollerCollector struct {
	Targets             []config.TargetConfig
	customLabelSet      *customLabelSet
	lookupTimes         map[string]map[string]float64
	availabilityStatus  map[string]bool
	dnsLookupTimeDesc   *prometheus.Desc
	dnsAvailabilityDesc *prometheus.Desc
}

func NewDNSPollerCollector(targets []config.TargetConfig) *DNSPollerCollector {
	return &DNSPollerCollector{
		Targets:            targets,
		customLabelSet:     newCustomLabelSet(targets),
		lookupTimes:        make(map[string]map[string]float64),
		availabilityStatus: make(map[string]bool),
		dnsLookupTimeDesc:   prometheus.NewDesc(dnsLookupTimeMetric.Name, dnsLookupTimeMetric.Help, append(newCustomLabelSet(targets).labelNames(), dnsLookupTimeMetric.LabelNames...), nil),
		dnsAvailabilityDesc: prometheus.NewDesc(dnsAvailabilityMetric.Name, dnsAvailabilityMetric.Help, append(newCustomLabelSet(targets).labelNames(), dnsAvailabilityMetric.LabelNames...), nil),
	}
}

func (collector *DNSPollerCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.dnsLookupTimeDesc
	ch <- collector.dnsAvailabilityDesc
}

func (collector *DNSPollerCollector) Collect(ch chan<- prometheus.Metric) {
	for _, target := range collector.Targets {
		labelValues := collector.customLabelSet.labelValues(target)

		if availability, exists := collector.availabilityStatus[target.Addr]; exists {
			ch <- prometheus.MustNewConstMetric(collector.dnsAvailabilityDesc, prometheus.GaugeValue, boolToFloat64(availability), append(labelValues, target.Addr)...)
		}

		if lookupTimesForServer, exists := collector.lookupTimes[target.Addr]; exists {
			for hostname, lookupTime := range lookupTimesForServer {
				ch <- prometheus.MustNewConstMetric(collector.dnsLookupTimeDesc, prometheus.GaugeValue, lookupTime, append(labelValues, target.Addr, hostname)...)
			}
		}
	}
}

func boolToFloat64(b bool) float64 {
    if b {
        return 1
    }
    return 0
}

func (collector *DNSPollerCollector) UpdateMetric(metricType, dnsServerLabel, hostnameLabel string, value float64) {
    switch metricType {
    case "dns_availability":
        collector.availabilityStatus[dnsServerLabel] = value == 1
    case "dns_lookup_time":
        if _, exists := collector.lookupTimes[dnsServerLabel]; !exists {
            collector.lookupTimes[dnsServerLabel] = make(map[string]float64)
        }
        collector.lookupTimes[dnsServerLabel][hostnameLabel] = value
    default:
        log.Errorln("UpdateMetric: Unknown metric name.")
    }
}