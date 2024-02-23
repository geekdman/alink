package ssh

import (
	"errors"
	"fmt"
	"github.com/pkg/sftp"
	"go.uber.org/zap"
	"golang.org/x/crypto/ssh"
	"strings"
	"time"
)

// ssh 客户端连接实体
type SSH struct {
	IP string
	Port int
	Username string
	Mode string   // 认证方式password  密码认证；key 密钥认证
	Password string
	Key string
	Client *ssh.Client
	SftpClient *sftp.Client
	LastResult string
}

func NewSSH(ip string, username string, password string, key string, mode string, port ...int) *SSH {
	/**
	  创建命令行实例
	  @param ip IP地址
	  @param username 用户名
	  @param password 登陆密码
	  @param key 认证私钥
	  @param mode 认证模式( password: 密码 | key: 秘钥 )
	  @param port 端口号, 不填写则默认为22
	*/

	client := new(SSH)
	client.IP = ip
	client.Username = username
	client.Password = password
	client.Key = key
	client.Mode = mode

	if len(port) <= 0 {
		client.Port = 22
	} else {
		client.Port = port[0]
	}

	return client
}

func (s *SSH) Connect() error {
	/**
	  SSH连接
	*/
	config := ssh.ClientConfig{
		User:            s.Username,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         10 * time.Second,
	}

	// 判断SSH连接的认证方式
	if s.Mode == "key" {
		signer, err := ssh.ParsePrivateKey([]byte(s.Key))
		if err != nil {
			zap.S().Fatalf("ssh key signer failed: %v", err)
		}
		config.Auth = []ssh.AuthMethod{ssh.PublicKeys(signer)}
	} else {
		fmt.Println("s.Password", s.Password)
		config.Auth = []ssh.AuthMethod{ssh.Password(s.Password)}
	}

	// 创建SSH客户端
	addr := fmt.Sprintf("%s:%d", s.IP, s.Port)
	sshClient, err := ssh.Dial("tcp", addr, &config)
	if err != nil {
		return err
	}
	s.Client = sshClient
	return nil
}

func (s SSH) Run(command string) (string, error) {
	/**
	  执行Shell命令
	  @param command 要执行的命令，多个命令采用 ; 隔开
	*/
	if s.Client == nil {
		if err := s.Connect(); err != nil {
			return "", err
		}
	}
	// 创建一个客户端SSH
	session, err := s.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	// 执行Shell命令
	buf, err := session.CombinedOutput(command)

	s.LastResult = string(buf)
	return s.LastResult, err
}

func (s SSH) AddPublicKeyToRemoteHost(publicKey string) error {
	/**
	  将公钥写入目标主机
	  700 是文档拥有可读可写可执行，同一组用户或者其他用户都不具有操作权限
	  600 是文件拥有者可读可写，不可执行，同一组用户或者其他用户都不具有操作权限
	*/
	command := fmt.Sprintf("mkdir -p -m 700 ~/.ssh && echo %v >> ~/.ssh/authorized_keys && chmod 600 ~/.ssh/authorized_keys", strings.TrimSpace(publicKey))
	_, err := s.Run(command)
	if err != nil {
		message := fmt.Sprintf("%v: %v", "添加公钥失败！", err)
		return errors.New(message)
	}
	return nil
}