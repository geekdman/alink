package zk

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/your/repo/utils/base"
	"github.com/samuel/go-zookeeper/zk"
)

type ZkCrud struct {
	ip      string
	port    string
	conn    *zk.Conn
}

func NewZkCrud(info string) *ZkCrud {
	parts := strings.Split(info, ":")
	if len(parts) != 2 {
		log.Fatal("Invalid ZooKeeper connection info. Use 'ip:port' format.")
	}

	return &ZkCrud{
		ip:   parts[0],
		port: parts[1],
	}
}

func (z *ZkCrud) connect() {
	address := fmt.Sprintf("%s:%s", z.ip, z.port)
	conn, _, err := zk.Connect([]string{address}, 10 * time.Second)
	if err != nil {
		log.Fatal("Error connecting to ZooKeeper:", err)
	}

	z.conn = conn
}

func (z *ZkCrud) Create(path, data string, ephemeral bool) (bool, interface{}) {
	if z.conn == nil {
		z.connect()
	}

	flags := 0
	if ephemeral {
		flags = zk.FlagEphemeral
	}

	path, err := z.conn.Create(path, []byte(data), flags, zk.WorldACL(zk.PermAll))
	if err != nil {
		return false, err
	}

	return true, fmt.Sprintf("Node created at %s", path)
}

func (z *ZkCrud) Read(path string) (bool, interface{}) {
	if z.conn == nil {
		z.connect()
	}

	data, _, err := z.conn.Get(path)
	if err != nil {
		return false, err
	}

	return true, string(data)
}

func (z *ZkCrud) Update(path, newData string) (bool, interface{}) {
	return z.Create(path, newData, false)
}

func (z *ZkCrud) Delete(path string) (bool, interface{}) {
	if z.conn == nil {
		z.connect()
	}

	err := z.conn.Delete(path, -1)
	if err != nil {
		return false, err
	}

	return true, "Node deleted successfully"
}

func (z *ZkCrud) Stats(query interface{}) (bool, interface{}) {
	// 实现统计操作
}

func (z *ZkCrud) Run() {
	z.connect()
	fmt.Printf("\n=== ZooKeeper Shell (Connected to %s:%s) ===\n", z.ip, z.port)
	fmt.Println("Commands: set, get, ls, stat, delete, help, exit")

	for {
		fmt.Print("zk> ")
		var cmd string
		fmt.Scanln(&cmd)

		cmd = strings.TrimSpace(cmd)
		if cmd == "" {
			continue
		}

		if cmd == "exit" {
			fmt.Println("Returning to main shell...")
			break
		}

		parts := strings.SplitN(cmd, " ", 2)
		if len(parts) < 1 {
			continue
		}

		action := parts[0]
		if action == "help" {
			fmt.Println("\nAvailable Commands:")
			fmt.Println("  set <path>=<data> [ephemeral] - Create or update a znode")
			fmt.Println("  get <path>                     - Retrieve data from a znode")
			fmt.Println("  ls [path]                      - List children of a znode")
			fmt.Println("  stat                           - Show ZooKeeper server statistics")
			fmt.Println("  delete <path>                  - Delete a znode")
			fmt.Println("  exit                           - Exit the ZooKeeper shell")
			fmt.Println("  help                           - Show this help message")
			continue
		}

		// 处理其他命令
	}
}