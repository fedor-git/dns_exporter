package dnspoller

import (
	"strconv"
	"time"

	"github.com/fedor-git/dns_exporter/core/config"
	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var Collector *DNSPollerCollector

func InitMetrics(conf *config.Config) {
	Collector = NewDNSPollerCollector(conf.Servers)
	prometheus.MustRegister(Collector)
}

func StartDNSThread(conf *config.Config) {
	go dnsTestThread(conf)
}

func checkDNSResolution(domain, dnsServer string, timeout time.Duration) float64 {
	resolver := dns.Client{Timeout: timeout}
	msg := new(dns.Msg)
	msg.SetQuestion(domain+".", dns.TypeA)

	startTime := time.Now()
	_, _, err := resolver.Exchange(msg, dnsServer+":53")
	endTime := time.Now()

	if err != nil {
		return 0.0
	}

	return endTime.Sub(startTime).Seconds()
}

func dnsTestThread(conf *config.Config) {
	for {
		type CheckResult struct {
			SuccessCount int
			TotalCount   int
		}

		serverChecks := make(map[string]*CheckResult)

		for _, t := range conf.Servers {
			serverChecks[t.Addr] = &CheckResult{TotalCount: len(conf.Hosts)}

			for _, d := range conf.Hosts {
				ck := checkDNSResolution(d, t.Addr, 1*time.Second)
				if ck == 0.0 {
					log.Debugf(`Lookup Time [%s]: %s = none`, t.Addr, d)
					Collector.UpdateMetric("dns_lookup_time", t.Addr, d, 0.0)
				} else {
					serverChecks[t.Addr].SuccessCount++
					log.Debugf(`Lookup Time [%s]: %s = %.4f`, t.Addr, d, ck)
					Collector.UpdateMetric("dns_lookup_time", t.Addr, d, ck)
				}
			}
		}
		for dnsServer, check := range serverChecks {
			log.Debugf(`Server: %s SuccessCount: %s TotalCount: %s`, dnsServer, strconv.Itoa(check.SuccessCount), strconv.Itoa(check.TotalCount))
			if check.SuccessCount > 0 {
				log.Debugf(`DNS Availability [%s]: true`, dnsServer)
				Collector.UpdateMetric("dns_availability", dnsServer, "", 1)
			} else {
				log.Debugf(`DNS Availability [%s]: false`, dnsServer)
				Collector.UpdateMetric("dns_availability", dnsServer, "", 0)
			}
		}

		if conf.Configuration.Write {
			log.Debug("Write to Files is Enabled!")
			fileResolvTime := conf.Configuration.Path + "/" + conf.Configuration.TimeFile
			fileDnsUpPath := conf.Configuration.Path + "/" + conf.Configuration.UPFile

			if err := gatherAndWriteMetrics([]string{"dns_availability"}, fileDnsUpPath); err != nil {
				log.Error("Error writing availability metrics:", err)
			}

			if err := gatherAndWriteMetrics([]string{"dns_lookup_time"}, fileResolvTime); err != nil {
				log.Error("Error writing lookup time metrics:", err)
			}
		}

		log.Debug("Job finished!")
		log.Debugf("Waiting new job: %s", time.Duration(conf.Configuration.Interval) * time.Second)
		time.Sleep(time.Duration(conf.Configuration.Interval) * time.Second)
	}
}
