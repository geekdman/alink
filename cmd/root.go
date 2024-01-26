/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"alink/config"
	"github.com/spf13/cobra"
	"os"
)

var (
	confpath string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "oeoslink",
	Short: "A brief description of your application",
	TraverseChildren: true,
	Long: `这是一个连接工具，用来连接zk、kafka、tikv、elasticsearch`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&confpath,"config","c","/usr/local/oct/oeos/conf/","-c /etc/oeos")
	config.GetConfig(confpath)
}


