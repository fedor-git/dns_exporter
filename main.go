package main

import (
	"fmt"
	"os"

	"github.com/fedor-git/dns_exporter/core/config"
	"github.com/fedor-git/dns_exporter/core/dnspoller"
	"github.com/fedor-git/dns_exporter/core/webapp"
	log "github.com/sirupsen/logrus"
	"gopkg.in/alecthomas/kingpin.v2"
)

const (
	internalPort        = 9433
	version      string = "0.0.1"
)

var (
	showVersion   = kingpin.Flag("version", "Print version information").Default().Bool()
	listenAddress = kingpin.Flag("web.listen-address", "Address on which to expose metrics and web interface").Default("").String()
	metricsPath   = kingpin.Flag("web.telemetry-path", "Path under which to expose metrics").Default("/metrics").String()
	configFile    = kingpin.Flag("config.path", "Path to config file").Default("").String()
	logLevel      = kingpin.Flag("log.level", "Only log messages with the given severity or above. Valid levels: [debug, info, warn, error, fatal]").Default("info").String()
)

func printVersion() {
	fmt.Println("dns-exporter")
	fmt.Printf("Version: %s\n", version)
	fmt.Println("Author(s): Fedir Sorokin")
	fmt.Println("Metric exporter for DNS infrastructure")
}

func loadConfig() (*config.Config, error) {
	if *configFile == "" {
		return nil, fmt.Errorf("config file path is empty")
	}

	f, err := os.Open(*configFile)
	if err != nil {
		return nil, fmt.Errorf("cannot load config file: %w", err)
	}
	defer f.Close()

	cfg, err := config.FromYAML(f)
	return cfg, err
}

func main() {
	kingpin.Parse()

	setLogLevel(*logLevel)
	log.SetReportCaller(true)

	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	cfg, err := loadConfig()
	if err != nil {
		kingpin.FatalUsage("could not load config.path: %v", err)
	}

	dnspoller.InitMetrics()
	dnspoller.StartDNSThread(cfg)

	webapp.StartServer(*listenAddress, internalPort, *metricsPath)
}
