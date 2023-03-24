package metrics

import (
	"fmt"
	"net/http"

	"github.com/Volkov-Stanislav/checker-reporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Instance struct {
	cfg               *config.Config
	checkPingTime     *prometheus.GaugeVec
	checkPingLost     *prometheus.GaugeVec
	checkHTTPTime     *prometheus.GaugeVec
	checkHTTPBandwith *prometheus.GaugeVec
}

func NewPrometheusInstance(cfg *config.Config) *Instance {
	var result Instance
	result.cfg = cfg
	result.register()

	return &result
}

func (o *Instance) Run() {
	go func() {
		http.Handle("/metrics", promhttp.Handler())
		err := http.ListenAndServe(":"+fmt.Sprint(o.cfg.MetricsPort), nil)
		
		if err != nil {
			panic(err)
		}
	}()
}

func (o *Instance) register() {
	// регистрирует дополнительные метрики.
	o.checkHTTPTime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "checker_reporter_check_http_time",
			Help: "Check http response time in seconds to endpoints.",
		},
		[]string{"sourceIP", "destIP"},
	)
	o.checkHTTPBandwith = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "checker_reporter_check_http_bandwidth",
			Help: "Check http response bandwidth in bytes/sec to endpoints.",
		},
		[]string{"sourceIP", "destIP"},
	)
	o.checkPingTime = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "checker_reporter_check_ping_time",
			Help: "Check ping average time in seconds to endpoints.",
		},
		[]string{"sourceIP", "destIP"},
	)
	o.checkPingLost = promauto.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "checker_reporter_check_ping_lost",
			Help: "Check ping lost packets to endpoints.",
		},
		[]string{"sourceIP", "destIP"},
	)
}

func (o *Instance) UpdateHTTP(time, bandwidth float64, destIP string) {
	o.checkHTTPBandwith.WithLabelValues(o.cfg.LocalIP, destIP).Set(bandwidth)
	o.checkHTTPTime.WithLabelValues(o.cfg.LocalIP, destIP).Set(time)
}

func (o *Instance) UpdatePING(time, lost float64, destIP string) {
	o.checkPingLost.WithLabelValues(o.cfg.LocalIP, destIP).Set(lost)
	o.checkPingTime.WithLabelValues(o.cfg.LocalIP, destIP).Set(time)
}
