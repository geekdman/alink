package main

import (
	"alink/pkg/ssh"
	"fmt"
	"os"
)

func main() {
    Host := "10.0.11.99"
    User := "root"
    Password := "password"
    Port := 22
    Key := ""
    Mode := "password"

    // 首次SSH连接使用密码
    cli := ssh.NewSSH(Host, User, Password, Key, Mode, Port)
    fmt.Println(cli.Password, cli.Username, cli.IP)
    if err := cli.Connect(); err != nil {
        print("连接失败！", err.Error())
        return
    } else {
        print("连接成功！")
    }

    publicKey, _ := os.ReadFile("./keys/id_rsa.pub") // 这里放公钥字符串
    fmt.Println("publicKey:::", string(publicKey))
    err := cli.AddPublicKeyToRemoteHost(string(publicKey)) // 将本地的公钥发送到远程服务器的指定位置
    if err != nil {
        panic(err)
    }

}