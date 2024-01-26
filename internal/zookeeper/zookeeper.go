package zookeeper

import (
	"alink/config"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"time"
)

type zkConn  struct {
	Conn *zk.Conn
	Path string
	Value string

}

func GetZKConn()  *zkConn{
	zC := new(zkConn)
	hosts := config.Cfg.GetZKConfig()
	conn, _, err := zk.Connect(hosts, time.Second*5)
	zC.Conn = conn
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//fmt.Println(conn.Server())
	return zC
}

// 新增节点
func (zc *zkConn) Add()  {

	exists,_,err := zc.Conn.Exists(zc.Path)
	if err != nil {
		log.Fatalf("查询%s失败, err: %v\n", zc.Path, err)
	}
	if exists {
		log.Fatalf("查询%s节点已经存在，请先删除改节点！！！", zc.Path)
	} else {
		acl := zk.WorldACL(zk.PermAll)
		var flags int32 = 0
		_, err := zc.Conn.Create(zc.Path,[]byte(zc.Value),flags,acl)
		if err != nil {
			log.Fatal(err)
		}
	}
}
// 查询node
func (zc *zkConn) Get() {

	res,s,err := zc.Conn.Children(zc.Path)

	if err != nil {
		log.Fatalf("查询%s失败, err: %v\n", zc.Path, err)
	}
	if s.NumChildren > 0 {
		fmt.Println("=========================")
		fmt.Printf("%s 存在子节点，如下:\n", zc.Path)
		fmt.Println(res)
		fmt.Println("=========================")
	} else if s.NumChildren == 0 {
		//data, _, err := conn.Get(path)
		//if err != nil {
		//	fmt.Printf("查询%s失败, err: %v\n", path, err)
		//	return
		//}
		//formatprint(data,path)
		data, _, err := zc.Conn.Get(zc.Path)
		if err != nil {
			fmt.Printf("查询%s失败, err: %v\n", zc.Path, err)
			return
		}
		fmt.Println("=========================")
		fmt.Printf("%s 的值为:\n", zc.Path)
		fmt.Println(string(data))
		fmt.Println("=========================")
	}
}
// 修改节点
func (zc *zkConn) Set() {
	// 获取 path 的属性
	_, stat, err := zc.Conn.Get(zc.Path)
	if err != nil {
		log.Fatal(err)
	}
	stat, err = zc.Conn.Set(zc.Path, []byte(zc.Value), stat.Version)
	if err != nil {
		log.Fatal(err)
	}
}
// 删除节点
func (zc *zkConn) Delete()  {
	_, stat, err :=  zc.Conn.Get(zc.Path)
	if err != nil {
		log.Fatal(err)
	}

	// version是用于 CAS支持，可以通过此种方式保证原子性
	if err := zc.Conn.Delete(zc.Path, stat.Version); err != nil {
		log.Fatal(err)
	}
}