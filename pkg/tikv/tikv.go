package tikv

import (
	"github.com/tikv/client-go/v2/txnkv"
)

type KV  struct {
	K,V []byte
}

type kvClient struct {
	 txn *txnkv.KVTxn
}

//func initStore()  {
//	var err error
//	client, err = txnkv.NewClient([]string{*pdAddr})
//	if err != nil {
//		panic(err)
//	}
//}
//
//func GetKVConn()  *kvClient{
//	hosts := config.Cfg.GetKVConfig()
//	txnkv.NewClient(hosts)
//	return nil
//}

//func main()  {
//
//}