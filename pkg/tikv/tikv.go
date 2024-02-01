//package tikv
//
//import (
//	"alink/config"
//	"fmt"
//	"time"
//	"github.com/pingcap/tidb/kv"
//	"github.com/pingcap/tidb/store/tikv"
//)
//
//type kvConn  struct {
//
//	K,V []byte
//}
//
//func GetKVConn()  *kvConn{
//	driver := tikv.Driver{}
//	var err error
//	store, err = driver.Open(fmt.Sprintf("tikv://%s", *pdAddr))
//	//fmt.Println(conn.Server())
//	return zC
//}