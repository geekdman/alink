/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"alink/internal/zookeeper"
	"fmt"
	"github.com/spf13/cobra"
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
		zc :=zookeeper.GetZKConn()
		// 关闭zk 连接
		defer zc.Conn.Close()
		zc.Path = path
		zc.Get()
	},
}

var addzkCmd = &cobra.Command{
	Use:   "add",
	Short: "add",
	Long: "add",
	Run: func(cmd *cobra.Command, args []string) {
		zc :=zookeeper.GetZKConn()
		// 关闭zk 连接
		defer zc.Conn.Close()
		zc.Path = path
		zc.Value = value
		zc.Add()
		zc.Get()
	},
}

var setzkCmd = &cobra.Command{
	Use:   "set",
	Short: "set",
	Long: "set",
	Run: func(cmd *cobra.Command, args []string) {
		zc :=zookeeper.GetZKConn()
		// 关闭zk 连接
		defer zc.Conn.Close()
		zc.Path = path
		zc.Value = value
		zc.Set()
		zc.Get()
	},
}

var delzkCmd = &cobra.Command{
	Use:   "del",
	Short: "del",
	Long: "del",
	Run: func(cmd *cobra.Command, args []string) {
		zc :=zookeeper.GetZKConn()
		// 关闭zk 连接
		defer zc.Conn.Close()
		zc.Path = path
		zc.Value = value
		zc.Delete()
		zc.Get()
	},
}

func init() {
	rootCmd.AddCommand(zkCmd)
	zkCmd.AddCommand(
		addzkCmd,
		delzkCmd,
		setzkCmd,
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
