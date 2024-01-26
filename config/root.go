package config

import (
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	V *viper.Viper
	Configpath  string
}

var Cfg *Config

func (cfg *Config)setConfig(confpath string)  *Config {
	viper.SetConfigName("system.config")
	viper.SetConfigType("properties")

	viper.AddConfigPath(confpath)
	viper.AddConfigPath("./")

	err := viper.ReadInConfig()
	if err != nil {
		log.Fatal("read config failed: %v", err)
	}
	cfg.V = viper.GetViper()
	return cfg
}

func GetConfig(confpath string)  *Config {
	c := new(Config)
	Cfg = c.setConfig(confpath)
	return Cfg
}

func init()  {
}

//func GetZookeeperServer()  {
//	getConfig()
//	v := Cfg.GetString("zookeeper.conf.client.servers")
//	fmt.Println(v)
//}