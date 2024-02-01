package config

import "strings"

func (cfg *Config) GetKVConfig()  []string {

	hosts := strings.Split(cfg.V.GetString("tikv.conf.pds.servers"),",")
	return hosts
}
