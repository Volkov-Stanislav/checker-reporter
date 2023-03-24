package checkers

import (
	"fmt"
	"sync"
	"time"

	"github.com/Volkov-Stanislav/checker-reporter/config"
	"github.com/Volkov-Stanislav/checker-reporter/metrics"
	probing "github.com/prometheus-community/pro-bing"
)

const (
	pingCount    = 4
	pingTimeout  = 10
	pingLostTime = 1000
)

type CheckPing struct {
	prom *metrics.Instance ``
	cfg  *config.Config
}

func NewCheckPing(prom *metrics.Instance, cfg *config.Config) *CheckPing {
	var result CheckPing
	result.prom = prom
	result.cfg = cfg
	return &result
}

func (o *CheckPing) Run(address string, wg *sync.WaitGroup) {
	defer wg.Done()

	pinger, err := probing.NewPinger(address)
	if err != nil {
		panic(err)
	}

	pinger.SetPrivileged(true)

	pinger.Count = pingCount
	pinger.Timeout = time.Duration(time.Second * pingTimeout)
	err = pinger.Run()

	if err != nil {
		panic(err)
	}

	var tim, lost float64

	stat := pinger.Statistics()
	for _, rtts := range stat.Rtts {
		tim += rtts.Seconds()
	}

	tim /= float64(pingCount)
	if tim == 0 {
		tim = pingLostTime
	}
	lost = stat.PacketLoss

	fmt.Printf("Ping %s: %#v \n", address, stat)
	o.prom.UpdatePING(tim, lost, address)
}
