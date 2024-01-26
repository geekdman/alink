package config

import "strings"

func (cfg *Config) GetZKConfig()  []string {

	hosts := strings.Split(cfg.V.GetString("zookeeper.conf.client.servers"),",")
	return hosts
}
