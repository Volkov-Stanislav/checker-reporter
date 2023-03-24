package main

import (
	"fmt"

	"github.com/Volkov-Stanislav/checker-reporter/checkers"
	"github.com/Volkov-Stanislav/checker-reporter/config"
	"github.com/Volkov-Stanislav/checker-reporter/httpoint"
	"github.com/Volkov-Stanislav/checker-reporter/metrics"
	"github.com/Volkov-Stanislav/checker-reporter/service"
)

const pingCount = 4

func main() {
	cf := config.NewConfig()
	//cf.Parse()
	cf.ConsulParse()
	fmt.Printf("Config: %#v \n", cf)

	pr := metrics.NewPrometheusInstance(cf)
	pr.Run()

	httpPoint := httpoint.NewHTTPoint(pr, cf.CheckPort)
	httpPoint.Run()

	pg := checkers.NewCheckPing(pr, cf)
	hck := checkers.NewHTTPCheck(pr, cf)

	r := service.NewRunner(cf, pr)
	r.Run(pg, hck)
}
