package config

import (
	"log"
	"strings"

	"github.com/spf13/viper"
)

var Cfg *viper.Viper

func setConfig() *viper.Viper {
	viper.SetConfigName("system.config")
	viper.SetConfigType("properties")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config failed: %v", err)
	}
	return viper.GetViper()
}

func GetConfig() {
	Cfg = setConfig()
}

func GetZKConfig(c *viper.Viper)  []string {
	hosts := strings.Split(c.GetString("zookeeper.conf.client.servers"),",")
	return hosts
}

func init()  {
	GetConfig()
}

//func GetZookeeperServer()  {
//	getConfig()
//	v := Cfg.GetString("zookeeper.conf.client.servers")
//	fmt.Println(v)
//}