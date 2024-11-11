/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"alink/pkg/ssh"
	"fmt"
	"github.com/spf13/cobra"
)

var (
	ip string
	username string
	password string
	key string
	mode string
	port int
	command string
)

// sshCmd represents the ssh command
var sshCmd = &cobra.Command{
	Use:   "ssh",
	Short: "connect server by ssh",
	Long: `connect server by ssh`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ssh called")
		var a = ssh.NewSSH(ip,username,password)
		a.Connect()
		fmt.Println(a.Run(command))
	},

}

func init() {
	rootCmd.AddCommand(sshCmd)

	sshCmd.Flags().StringVarP(&ip,"ip","H","","eg: http://localhost:9000 ")
	sshCmd.Flags().StringVarP(&username,"username","u","","")
	sshCmd.Flags().StringVarP(&password,"password","p","","")
	sshCmd.Flags().StringVarP(&key,"privatekey","i","","")
	sshCmd.Flags().StringVarP(&mode,"mode","m","password","default is password")
	sshCmd.Flags().IntVarP(&port,"port","P",22,"default is 22")
	sshCmd.Flags().StringVarP(&command,"command","c","","使用双引号包裹命令")
	//sshCmd.Flags().StringSliceVarP(&command,"command","c",nil,"default")
}
