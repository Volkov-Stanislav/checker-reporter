package httpoint

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/Volkov-Stanislav/checker-reporter/metrics"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type HTTPoint struct {
	prom   *metrics.Instance
	buffer []byte
	port   int
}

const replayHTTPcheckSize = 1024 * 1024

func NewHTTPoint(prom *metrics.Instance, port int) *HTTPoint {
	var result HTTPoint
	result.prom = prom
	result.port = port
	result.setBuffer(replayHTTPcheckSize)

	return &result
}

func (o *HTTPoint) Run() {
	go func() {
		err := o.serve()
		if err != nil {
			panic(err)
		}
	}()
}

func (o *HTTPoint) serve() error {
	http.Handle("/reporter", promhttp.InstrumentHandlerCounter(
		promauto.NewCounterVec(
			prometheus.CounterOpts{
				Name: "checker_reporter_requests_httpoint_total",
				Help: "Total number of httpoint requests by HTTP code.",
			},
			[]string{"code"},
		),
		http.HandlerFunc(o.getReporter),
	))

	return http.ListenAndServe(":"+fmt.Sprintf("%d", o.port), nil)
}

func (o *HTTPoint) getReporter(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("got /reporter request\n")
	
	_, err := w.Write(o.buffer)
	if err != nil {
		fmt.Printf("replay error \n")
	}
}

func (o *HTTPoint) setBuffer(n int) {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Int63()%int64(len(letterBytes))]
	}

	o.buffer = b
}
