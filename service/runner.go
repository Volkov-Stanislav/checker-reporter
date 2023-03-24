package service

import (
	"fmt"
	"sync"
	"time"

	"github.com/Volkov-Stanislav/checker-reporter/config"
	"github.com/Volkov-Stanislav/checker-reporter/metrics"
)

type Checker interface {
	Run(address string, wg *sync.WaitGroup)
}

type Runner struct {
	cfg  *config.Config
	prom *metrics.Instance
}

func NewRunner(cfg *config.Config, prom *metrics.Instance) *Runner {
	var result Runner
	result.cfg = cfg
	result.prom = prom

	return &result
}

func (o *Runner) Run(checks ...Checker) {
	var wg sync.WaitGroup

	for {
		for ip := range o.cfg.HostsMap {
			for _, chk := range checks {
				wg.Add(1)
				fmt.Printf("RunCheck ip: %s,  chk: %#v \n", ip, chk)

				go chk.Run(ip, &wg)
			}
		}

		wg.Wait()
		time.Sleep(time.Duration(o.cfg.UpdateInterval) * time.Second)
	}
}
