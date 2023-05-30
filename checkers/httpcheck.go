package checkers

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/Volkov-Stanislav/checker-reporter/config"
	"github.com/Volkov-Stanislav/checker-reporter/metrics"
)

type HTTPcheck struct {
	prom *metrics.Instance ``
	cfg  *config.Config
}

func NewHTTPCheck(prom *metrics.Instance, cfg *config.Config) *HTTPcheck {
	var result HTTPcheck
	result.prom = prom
	result.cfg = cfg

	return &result
}

func (o *HTTPcheck) Run(address string, wg *sync.WaitGroup) {
	defer wg.Done()

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	t := time.Now()

	url := "http://" + address + ":" + fmt.Sprint(o.cfg.CheckPort) + "/reporter"
	response, err := client.Get(url)

	if err != nil {
		fmt.Println(err)
		return
	}

	httpReplyTime := time.Since(t).Seconds()
	fmt.Printf("Reply from URL: %s in msec: %v\n", address, httpReplyTime)

	t = time.Now()
	// read response body
	_, err = io.ReadAll(response.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	// close response body
	response.Body.Close()

	httpReadTime := time.Since(t).Seconds()
	httpBandwith := 1048576 / httpReadTime

	fmt.Printf("Read URL: %s  in sec: %v  bandwidth %v \n", address, httpReadTime, httpBandwith)
	o.prom.UpdateHTTP(httpReplyTime, httpBandwith, address)
}
