/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"alink/pkg/s3"
	"fmt"

	"github.com/spf13/cobra"
)


var (
	endpoint string
	accessKeyID string
	secretAccessKey string
	useSSL bool
	bucketname string
)
// s3clientCmd represents the s3client command
var s3clientCmd = &cobra.Command{
	Use:   "s3",
	Short: "A brief description of your command",
	Long: `Manage s3 buckets and objects`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("s3client called")
	},
}

var copyCmd = &cobra.Command{
	Use:   "cp",
	Short: "copy objects",
	Long: "copy objects",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var listCmd = &cobra.Command{
	Use:   "ls",
	Short: "list buckets and objects",
	Long: "list buckets and objects",
	Run: func(cmd *cobra.Command, args []string) {
		s := s3.NewS3Conn(endpoint,accessKeyID,secretAccessKey,useSSL)
		s.List()
	},
}

var makebucketCmd = &cobra.Command{
	Use:   "mb",
	Short: "make a bucket",
	Long: "make a bucket",
	Run: func(cmd *cobra.Command, args []string) {
		s := s3.NewS3Conn(endpoint,accessKeyID,secretAccessKey,useSSL)
		s.MakeBucket(bucketname)
	},
}

var removebucketCmd = &cobra.Command{
	Use:   "rb",
	Short: "remove a bucket",
	Long: "remove a bucket",
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.AddCommand(s3clientCmd)

	s3clientCmd.AddCommand(
		copyCmd,
		listCmd,
		makebucketCmd,
		removebucketCmd,
		)

	copyCmd.Flags().StringVarP(&endpoint,"url","u","","eg: http://localhost:9000 ")
	copyCmd.Flags().StringVarP(&accessKeyID,"access-key","k","","")
	copyCmd.Flags().StringVarP(&secretAccessKey,"secret-key","s","","")
	copyCmd.Flags().StringVarP(&bucketname,"bucketname","n","","")
	copyCmd.Flags().BoolVarP(&useSSL,"useSSL","",false,"default is false")

	listCmd.Flags().StringVarP(&endpoint,"url","u","","eg: http://localhost:9000 ")
	listCmd.Flags().StringVarP(&accessKeyID,"access-key","k","","")
	listCmd.Flags().StringVarP(&secretAccessKey,"secret-key","s","","")
	listCmd.Flags().StringVarP(&bucketname,"bucketname","n","","")
	listCmd.Flags().BoolVarP(&useSSL,"useSSL","",false,"default is false")

	makebucketCmd.Flags().StringVarP(&endpoint,"url","u","","eg: http://localhost:9000 ")
	makebucketCmd.Flags().StringVarP(&accessKeyID,"access-key","k","","")
	makebucketCmd.Flags().StringVarP(&secretAccessKey,"secret-key","s","","")
	makebucketCmd.Flags().StringVarP(&bucketname,"bucketname","n","","")
	makebucketCmd.Flags().BoolVarP(&useSSL,"useSSL","",false,"default is false")

	removebucketCmd.Flags().StringVarP(&endpoint,"url","u","","eg: http://localhost:9000 ")
	removebucketCmd.Flags().StringVarP(&accessKeyID,"access-key","k","","")
	removebucketCmd.Flags().StringVarP(&secretAccessKey,"secret-key","s","","")
	removebucketCmd.Flags().StringVarP(&bucketname,"bucketname","n","","")
	removebucketCmd.Flags().BoolVarP(&useSSL,"useSSL","",false,"default is false")
}
