package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gocolly/colly"

	"github.com/ghodss/yaml"
	_ "github.com/go-sql-driver/mysql"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/version"
	log "github.com/sirupsen/logrus"
)

var config Config

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
	flag.StringVar(&bind, "bind", "0.0.0.0:9104", "bind")
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
	if err := yaml.Unmarshal(b, &config); err != nil {
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
type Config struct {
	DSN     string
	Metrics map[string]struct {
		Query       string
		Type        string
		Description string
		Labels      []string
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
	for metricName, metric := range config.Metrics {
		metric.metricDesc = prometheus.NewDesc(
			prometheus.BuildFQName(collector, "", metricName),
			metric.Description,
			metric.Labels, nil,
		)
		config.Metrics[metricName] = metric
		log.Infof("metric description for \"%s\" registerd", metricName)
	}
}

// Collect prometheus collect
func (e *QueryCollector) Collect(ch chan<- prometheus.Metric) {

	c := colly.NewCollector()
	c.OnHTML("div.official-name", func(e *colly.HTMLElement) {
		e.ForEach("span", func(_ int, el *colly.HTMLElement) {
			fmt.Println(el.ChildText("span:nth-child(1)"))
			/*writer.Write([]string{
				el.ChildText("td:nth-child(1)"),
				el.ChildText("td:nth-child(2)"),
				el.ChildText("td:nth-child(3)"),
			})*/
		})
		fmt.Println("Scrapping Complete")
	})
	//c.Visit("https://www.livecoinwatch.com/")
	c.Visit("https://www.livecoinwatch.com/price/Cardano-ADA")
	// Execute each queries in metrics
	//for name, metric := range config.Metrics {

	/*// Execute query
	rows, err := db.Query(metric.Query)
	if err != nil {
		log.Errorf("Failed to execute query: %s", err)
		continue
	}

	// Get column info
	cols, err := rows.Columns()
	if err != nil {
		log.Errorf("Failed to get column meta: %s", err)
		continue
	}

	des := make([]interface{}, len(cols))
	res := make([][]byte, len(cols))
	for i := range cols {
		des[i] = &res[i]
	}*/

	// fetch database
	/*for rows.Next() {
	rows.Scan(des...)
	data := make(map[string]string)
	for i, bytes := range res {
		data[cols[i]] = string(bytes)
	}

	// Metric labels
	labelVals := []string{}
	for _, label := range metric.Labels {
		labelVals = append(labelVals, data[label])
		//Replace labelVals with the string arrays
	}*/

	// Metric value
	//val, _ := strconv.ParseFloat(data[metric.Value], 64)
	//Replace val with the counter value

	// Add metric
	//ch <- prometheus.MustNewConstMetric(metric.metricDesc, prometheus.GaugeValue, val, labelVals...)

	//}
	//}
}
