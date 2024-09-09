package tikv

import (
	"alink/config"
	"fmt"
	"time"
    "github.com/tikv/client-go/v2/txnkv"

)

type kvConn  struct {
	K,V []byte
}

func GetKVConn()  *kvConn{
	hosts := config.Cfg.GetKVConfig()
	client, err := txnkv.NewClient(hosts)
	return nil
}