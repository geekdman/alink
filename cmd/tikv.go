/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// tikvCmd represents the tikv command
var tikvCmd = &cobra.Command{
	Use:   "tikv",
	Short: "connect tikv",
	Long: `connect tikv`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tikv called")
	},
}

func init() {
	rootCmd.AddCommand(tikvCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// tikvCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// tikvCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
