package dnspoller

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fedor-git/dns_exporter/core/config"
	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	dnsTimeMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cvpn_dns_resolution_time",
			Help: "Time of resolution",
		},
		[]string{"dns_server", "hostname"},
	)

	dnsUpMetric = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "cvpn_dns_up",
			Help: "Up state for DNS servers",
		},
		[]string{"dns_server"},
	)
)

func InitMetrics() {
	prometheus.MustRegister(dnsTimeMetric)
	prometheus.MustRegister(dnsUpMetric)
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
		serverStatus := make(map[string]bool)
		var tempResolv, tempUp []string

		tempResolv = append(tempResolv, "# HELP dns_lookup_time DNS lookup time measurement")
		tempResolv = append(tempResolv, "# TYPE dns_lookup_time gauge")

		tempUp = append(tempUp, "# HELP dns_availability DNS Server Status")
		tempUp = append(tempUp, "# TYPE dns_availability gauge")

		for _, t := range conf.Servers {
			for _, d := range conf.Hosts {
				ck := checkDNSResolution(d, t, 1*time.Second)

				if ck == 0.0 {
					tempResolv = append(tempResolv, fmt.Sprintf(`dns_lookup_time{dns_server="%s", hostname="%s"} 0.0`, t, d))
					dnsTimeMetric.WithLabelValues(t, d).Set(0)
					serverStatus[t] = false
				} else {
					tempResolv = append(tempResolv, fmt.Sprintf(`dns_lookup_time{dns_server="%s", hostname="%s"} %.4f`, t, d, ck))
					dnsTimeMetric.WithLabelValues(t, d).Set(ck)
					serverStatus[t] = true
				}
			}
		}
		if conf.Configuration.Write {

			for dnsServer, status := range serverStatus {
				if status {
					tempUp = append(tempUp, fmt.Sprintf(`dns_availability{dns_server="%s"} 1.0`, dnsServer))
					dnsUpMetric.WithLabelValues(dnsServer).Set(1)
				} else {
					tempUp = append(tempUp, fmt.Sprintf(`dns_availability{dns_server="%s"} 0.0`, dnsServer))
					dnsUpMetric.WithLabelValues(dnsServer).Set(0)
				}
			}

			fileResolvTime := conf.Configuration.Path + "/" + conf.Configuration.TimeFile
			fileDnsUpPath := conf.Configuration.Path + "/" + conf.Configuration.UPFile

			err := writeToFile(fileResolvTime, tempResolv)
			if err != nil {
				fmt.Println("Error writing to file:", err)
			}

			tempUp = uniqueStrings(tempUp)
			err = writeToFile(fileDnsUpPath, tempUp)
			if err != nil {
				fmt.Println("Error writing to file:", err)
			}
		}
		time.Sleep(time.Duration(conf.Configuration.Interval) * time.Second)
	}
}

func writeToFile(filePath string, lines []string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	content := strings.Join(lines, "\n") + "\n"
	_, err = file.WriteString(content)
	return err
}

func uniqueStrings(slice []string) []string {
	encountered := map[string]bool{}
	result := []string{}

	for _, v := range slice {
		if !encountered[v] {
			encountered[v] = true
			result = append(result, v)
		}
	}

	return result
}
