/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// esCmd represents the es command
var esCmd = &cobra.Command{
	Use:   "es",
	Short: "connect elaticsearch ",
	Long: `connect elasticsearch`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("es called")
	},
}

func init() {
	rootCmd.AddCommand(esCmd)

 	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// esCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// esCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
