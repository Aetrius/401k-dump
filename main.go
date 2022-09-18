package main

import (
	"flag"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/ghodss/yaml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

var coinConfig CoinConfig

const (
	collector = "query_exporter"
)

func main() {
	var err error
	var configFile, bind string
	// =====================
	// Get OS parameter
	// =====================
	flag.StringVar(&configFile, "config", "config.yml", "configuration file")
	flag.StringVar(&bind, "bind", "0.0.0.0:9112", "bind")
	flag.Parse()

	// =====================
	// Load config & yaml
	// =====================
	var b []byte
	if b, err = ioutil.ReadFile(configFile); err != nil {
		log.Errorf("Failed to read config file: %s", err)
		os.Exit(1)
	}

	// Load yaml
	if err := yaml.Unmarshal(b, &coinConfig); err != nil {
		log.Errorf("Failed to load config: %s", err)
		os.Exit(1)
	}

	// ========================
	// Regist handler
	// ========================
	log.Infof("Regist version collector - %s", collector)
	prometheus.Register(version.NewCollector(collector))
	prometheus.Register(&QueryCollector{})

	// Regist http handler
	log.Infof("HTTP handler path - %s", "/metrics")
	http.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		h := promhttp.HandlerFor(prometheus.Gatherers{
			prometheus.DefaultGatherer,
		}, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	})

	// start server
	log.Infof("Starting http server - %s", bind)
	if err := http.ListenAndServe(bind, nil); err != nil {
		log.Errorf("Failed to start http server: %s", err)
	}
}

// =============================
// Config config structure
// =============================
type CoinConfig struct {
	Metrics map[string]struct {
		URL         []string
		Qty         string
		Type        string
		Description string
		Value       string
		metricDesc  *prometheus.Desc
	}
}

// =============================
// QueryCollector exporter
// =============================
type QueryCollector struct{}

// Describe prometheus describe
func (e *QueryCollector) Describe(ch chan<- *prometheus.Desc) {
	for metricName, metric := range coinConfig.Metrics {
		name := "401kdb"
		metric.metricDesc = prometheus.NewDesc(
			prometheus.BuildFQName(collector, "", name),
			metric.Description,
			[]string{"retirement"}, nil,
		)
		coinConfig.Metrics[metricName] = metric
		log.Infof("metric description for \"%s\" registerd", metricName)
	}
}

// Collect prometheus collect
func (e *QueryCollector) Collect(ch chan<- prometheus.Metric) {
	log.Infof("metric description for \"%s\" registerd", "401kdb")
	for metricName, metric := range coinConfig.Metrics {

		result, err := strconv.ParseFloat(metric.Qty, 8)
		if err != nil {
			panic(err)
		}
		ch <- prometheus.MustNewConstMetric(metric.metricDesc, prometheus.GaugeValue, result, metricName)
	}
}
