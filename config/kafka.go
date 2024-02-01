package config

import "strings"

func (cfg *Config) GetKafkaConfig()  []string {

	hosts := strings.Split(cfg.V.GetString("kafka.conf.bootstrap.servers"),",")
	return hosts
}
