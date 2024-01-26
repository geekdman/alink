/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"alink/config"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/spf13/cobra"
	"log"
	"time"
)

var (
	path string
	value string
	filename string
)
// zookeeperCmd represents the zookeeper command
var zkCmd = &cobra.Command{
	Use:   "zk",
	Short: "zookeeper",
	Long: `Used to connect to zk to perform operations of adding, deleting, modifying and checking`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("zookeeper is connecting ")
		//getZKConn()
	},
}

var getzkCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long: "Get the value of the node. If there are child nodes, the child nodes will be returned.",
	Run: func(cmd *cobra.Command, args []string) {
		conn := getZKConn()
		// 关闭zk 连接
		defer conn.Close()
		get(conn,path)
	},
}

var addzkCmd = &cobra.Command{
	Use:   "add",
	Short: "add",
	Long: "add",
	Run: func(cmd *cobra.Command, args []string) {
		conn := getZKConn()
		// 关闭zk 连接
		defer conn.Close()
		add(conn,path,[]byte(value))
		get(conn,path)
	},
}

var setzkCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	Long: "set",
	Run: func(cmd *cobra.Command, args []string) {
		conn := getZKConn()
		// 关闭zk 连接
		defer conn.Close()
		set(conn,path,[]byte(value))
		get(conn,path)
	},
}

var delzkCmd = &cobra.Command{
	Use:   "del",
	Short: "del",
	Long: "del",
	Run: func(cmd *cobra.Command, args []string) {
		conn := getZKConn()
		// 关闭zk 连接
		defer conn.Close()
		delete(conn,path)
		get(conn,path)
	},
}

func init() {
	rootCmd.AddCommand(zkCmd)
	zkCmd.AddCommand(
		addzkCmd,
		delzkCmd,
		modifyzkCmd,
		getzkCmd,
		)
	//add 增加
	addzkCmd.Flags().StringVarP(&path,"key","k","","eg: /test ")
	addzkCmd.Flags().StringVarP(&value,"value","v","","")
	addzkCmd.Flags().StringVarP(&filename,"filename","f","","")
	//set 修改
	setzkCmd.Flags().StringVarP(&path,"key","k","","eg: /test")
	setzkCmd.Flags().StringVarP(&value,"value","v","","")
	setzkCmd.Flags().StringVarP(&filename,"filename","f","","")
	//
	//
	getzkCmd.Flags().StringVarP(&path,"key","k","","eg: /test")
	delzkCmd.Flags().StringVarP(&path,"key","k","","eg: /test")

}

func getZKConn()  *zk.Conn{
	hosts := config.Cfg.GetZKConfig()
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//fmt.Println(conn.Server())
	return conn
}

// 新增节点
func add(conn *zk.Conn, path string,data []byte)  {

	exists,_,err := conn.Exists(path)
	if err != nil {
		log.Fatalf("查询%s失败, err: %v\n", path, err)
	}
	if exists {
		log.Fatalf("查询%s节点已经存在，请先删除改节点！！！", path)
	} else {
		acl := zk.WorldACL(zk.PermAll)
		var flags int32 = 0
		_, err := conn.Create(path,data,flags,acl)
		if err != nil {
			log.Fatal(err)
		}
	}
}
// 查询node
func get(conn *zk.Conn,path string) {

	res,s,err := conn.Children(path)

	if err != nil {
		log.Fatalf("查询%s失败, err: %v\n", path, err)
	}
	if s.NumChildren > 0 {
		fmt.Println("=========================")
		fmt.Printf("%s 存在子节点，如下:\n", path)
		fmt.Println(res)
		fmt.Println("=========================")
	} else if s.NumChildren == 0 {
		//data, _, err := conn.Get(path)
		//if err != nil {
		//	fmt.Printf("查询%s失败, err: %v\n", path, err)
		//	return
		//}
		//formatprint(data,path)
		data, _, err := conn.Get(path)
		if err != nil {
			fmt.Printf("查询%s失败, err: %v\n", path, err)
			return
		}
		fmt.Println("=========================")
		fmt.Printf("%s 的值为:\n", path)
		fmt.Println(string(data))
		fmt.Println("=========================")
	}
}
// 修改节点
func set(conn *zk.Conn,path string,data []byte ) {
	// 获取 path 的属性
	_, stat, err := conn.Get(path)
	if err != nil {
		log.Fatal(err)
	}
	stat, err = conn.Set(path, data, stat.Version)
	if err != nil {
		log.Fatal(err)
	}
}
// 删除节点
func delete(conn *zk.Conn,path string)  {
	_, stat, err := conn.Get(path)
	if err != nil {
		log.Fatal(err)
	}

	// version是用于 CAS支持，可以通过此种方式保证原子性
	if err := conn.Delete(path, stat.Version); err != nil {
		log.Fatal(err)
	}
}