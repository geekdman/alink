/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/spf13/cobra"
	"alink/config"
	"time"
)

// zookeeperCmd represents the zookeeper command
var zkCmd = &cobra.Command{
	Use:   "zookeeper",
	Short: "zookeeper",
	Long: `zookeeper`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("zookeeper is connecting ")
		//getZKConn()
	},
}

var getzkCmd = &cobra.Command{
	Use:   "get",
	Short: "get",
	Long: "get",
	Run: func(cmd *cobra.Command, args []string) {
		get(path)
	},
}

var addzkCmd = &cobra.Command{
	Use:   "add",
	Short: "add",
	Long: "add",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("zookeeper add called")
	},
}

var modifyzkCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	Long: "set",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("zookeeper modify called")
	},
}

var delzkCmd = &cobra.Command{
	Use:   "del",
	Short: "del",
	Long: "del",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("zookeeper del called")
	},
}

var (
	path string
	value string
	filename string
)

func init() {
	rootCmd.AddCommand(zkCmd)
	zkCmd.AddCommand(
		addzkCmd,
		delzkCmd,
		modifyzkCmd,
		getzkCmd,
		)
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//zookeeperCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// zookeeperCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	addzkCmd.Flags().StringVarP(&path,"key","k","","eg: /test ")
	addzkCmd.Flags().StringVarP(&value,"value","v","","")
	addzkCmd.Flags().StringVarP(&filename,"filename","f","","")
	//
	modifyzkCmd.Flags().StringVarP(&path,"key","k","","eg: /test")
	modifyzkCmd.Flags().StringVarP(&value,"value","v","","")
	modifyzkCmd.Flags().StringVarP(&filename,"filename","f","","")
	//
	//
	getzkCmd.Flags().StringVarP(&path,"key","k","","eg: /test")
	delzkCmd.Flags().StringVarP(&path,"key","k","","eg: /test")

}

func getZKConn()  *zk.Conn{
	hosts := config.GetZKConfig(config.Cfg)
	conn, _, err := zk.Connect(hosts, time.Second*5)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	//fmt.Println(conn.Server())
	return conn
}

func ensure()  {

}

// 查询节点
func get(path string)  {
	conn := getZKConn()
	defer conn.Close()
	res,s,err := conn.Children(path)
	if err != nil {
		fmt.Printf("查询%s失败, err: %v\n", path, err)
		return
	}
	if s.NumChildren > 0 {
		fmt.Printf("%s 存在子节点，如下:\n", path)
		fmt.Println(res)
	} else {
		data, _, err := conn.Get(path)

		if err != nil {
			fmt.Printf("查询%s失败, err: %v\n", path, err)
			return
		}
		var prettyJSON bytes.Buffer
		error := json.Indent(&prettyJSON, data, "", "\t")
		if error != nil {
			fmt.Println(error)
			return
		}
		fmt.Printf("%s 的值为:\n", path)
		fmt.Println(string(prettyJSON.Bytes()))

	}
}
// 新增node
func add() {

}
// 修改节点
func modify()  {
	
}
// 删除节点
func delete()  {
	
}