package config

/*
New Config - array of

target_ip - ip-адрес машины на который установлено прила
metric_port - static port
hostname - имя хоста (согласовать маску для имен машин)
geo_loc - принадлежность к ЦОД по геолокации
virtual_type - тип виртуализации
location_type - {DC/MAG} - влияет паттерн проверок
*/

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/Volkov-Stanislav/checker-reporter/utils"
	capi "github.com/hashicorp/consul/api"
	"github.com/namsral/flag"
)

const (
	consulURL = "consul.service.datacenter_msk1.consul:8500"
)

type Config struct {
	UpdateInterval int    `json:"update_interval"`
	CheckHosts     string `json:"check_hosts"`
	MetricsPort    int    `json:"metrics_port"`
	CheckPort      int    `json:"check_port"`
	HostsMap       map[string]bool
	LocalIP        string
}

func NewConfig() *Config {
	var result Config
	result.HostsMap = make(map[string]bool)

	return &result
}

func (o *Config) Parse() {
	flag.String(flag.DefaultConfigFlagname, "config.conf", "path to config file")
	flag.IntVar(&o.UpdateInterval, "update_interval", 60, "scan interval")
	flag.StringVar(&o.CheckHosts, "check_hosts", "[127.0.0.1,10.0.0.0]", "array of all host checked, include local")
	flag.IntVar(&o.MetricsPort, "metrics_port", 2112, "local port for binding metrics handler")
	flag.IntVar(&o.CheckPort, "check_port", 8080, "local port for binding http check handler")
	flag.Parse()

	o.confProcess()
}

func (o *Config) ConsulParse() {
	fmt.Printf("BeginConnectToConsul\n")

	conf := &capi.Config{
		Address: consulURL,
	}
	client, err := capi.NewClient(conf)

	if err != nil {
		fmt.Printf("Error connect to Consul: %v", err)
		panic(err)
	}

	// Get a handle to the KV API
	kv := client.KV()

	// Lookup the pair
	pair := getValue(kv, "CheckerReporter/CheckHosts")
	if pair == nil {
		o.Parse()
		return
	}
	o.CheckHosts = string(pair.Value)
	fmt.Printf("KV: %v %s\n", pair.Key, pair.Value)

	pair = getValue(kv, "CheckerReporter/CheckPort")
	if pair == nil {
		o.Parse()
		return
	}
	o.CheckPort, _ = strconv.Atoi(string(pair.Value))
	fmt.Printf("KV: %v %s\n", pair.Key, pair.Value)

	pair = getValue(kv, "CheckerReporter/MetricsPort")
	if pair == nil {
		o.Parse()
		return
	}
	o.MetricsPort, _ = strconv.Atoi(string(pair.Value))
	fmt.Printf("KV: %v %s\n", pair.Key, pair.Value)

	pair = getValue(kv, "CheckerReporter/UpdateInterval")
	if pair == nil {
		o.Parse()
		return
	}
	o.UpdateInterval, _ = strconv.Atoi(string(pair.Value))
	fmt.Printf("KV: %v %s\n", pair.Key, pair.Value)

	o.confProcess()
}

func (o *Config) confProcess() {
	hosts := []string{}

	err := json.Unmarshal([]byte(o.CheckHosts), &hosts)
	if err != nil {
		return
	}

	for _, h := range hosts {
		o.HostsMap[h] = true
	}

	ips := utils.GetLocalIPAdresses()

	for _, ip := range ips {
		if _, ok := o.HostsMap[ip]; ok {
			o.LocalIP = ip
			delete(o.HostsMap, ip)

			break
		}
	}
}

func getValue(kv *capi.KV, name string) *capi.KVPair {
	pair, _, err := kv.Get(name, nil)
	if err != nil {
		fmt.Printf("Error Get KeyValue from Consul: %v", err)
		return nil
	}

	if pair == nil {
		fmt.Printf("Error Get KeyValue '%s' from Consul, not exist, fill from config file\n", name)
		return nil
	}

	return pair
}
