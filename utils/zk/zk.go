package zk

import (
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chzyer/readline"
	"github.com/samuel/go-zookeeper/zk"
)

type ZkCrud struct {
	ip       string
	port     string
	conn     *zk.Conn
	basePath string
}

func NewZkCrud(info string) *ZkCrud {
	parts := strings.Split(info, ":")
	if len(parts) != 2 {
		log.Fatal("Invalid ZooKeeper connection info. Use 'ip:port' format.")
	}

	return &ZkCrud{
		ip:       parts[0],
		port:     parts[1],
		basePath: "/",
	}
}

func (z *ZkCrud) connect() {
	address := fmt.Sprintf("%s:%s", z.ip, z.port)
	conn, _, err := zk.Connect([]string{address}, 10*time.Second)
	if err != nil {
		log.Fatal("Error connecting to ZooKeeper:", err)
	}

	z.conn = conn
}

func (z *ZkCrud) Create(data interface{}, identifier string) (bool, interface{}) {
	// 验证路径是否合法
	if !strings.HasPrefix(identifier, "/") {
		return false, fmt.Errorf("invalid path: %s", identifier)
	}

	if z.conn == nil {
		z.connect()
	}

	acl := zk.WorldACL(zk.PermAll)

	_, err := z.conn.Create(identifier, []byte(fmt.Sprint(data)), 0, acl)
	if err != nil {
		return false, err
	}

	return true, fmt.Sprintf("Node created at %s", identifier)
}

func (z *ZkCrud) Read(query interface{}) (bool, interface{}) {
	if z.conn == nil {
		z.connect()
	}

	path := query.(string)
	data, _, err := z.conn.Get(path)
	if err != nil {
		return false, err
	}

	return true, string(data)
}

func (z *ZkCrud) Update(identifier string, newData interface{}) (bool, interface{}) {
	return z.Create(newData, identifier)
}

func (z *ZkCrud) Delete(identifier string) (bool, interface{}) {
	if z.conn == nil {
		z.connect()
	}

	err := z.conn.Delete(identifier, -1)
	if err != nil {
		return false, err
	}

	return true, "Node deleted successfully"
}

func (z *ZkCrud) Stats(query interface{}) (bool, interface{}) {
	return true, "Statistics information is not directly available in the current ZooKeeper client version."
}

func (z *ZkCrud) Run() {
	z.connect()
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          "zk> ",
		HistoryFile:     "/tmp/zk_readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Printf("Error initializing readline: %v\n", err)
		return
	}
	defer rl.Close()

	fmt.Printf("\n=== ZooKeeper Shell (Connected to %s:%s) ===\n", z.ip, z.port)
	fmt.Println("Commands: set, get, ls, stat, delete, help, exit")

	for {
		line, err := rl.Readline()
		if err != nil {
			break
		}

		cmd := strings.TrimSpace(line)
		if cmd == "" {
			continue
		}

		parts := strings.SplitN(cmd, " ", 2)
		if len(parts) < 1 {
			continue
		}

		action := parts[0]
		switch action {
		case "exit":
			fmt.Println("Returning to main shell...")
			return
		case "help":
			fmt.Println("\nAvailable Commands:")
			fmt.Println("  set <path>=<data> - Create or update a znode")
			fmt.Println("  get <path>         - Retrieve data from a znode")
			fmt.Println("  ls [path]          - List children of a znode")
			fmt.Println("  stat               - Show ZooKeeper server statistics")
			fmt.Println("  delete <path>      - Delete a znode")
			fmt.Println("  exit               - Exit the ZooKeeper shell")
			fmt.Println("  help               - Show this help message")
		case "set":
			if len(parts) < 2 {
				fmt.Println("Usage: set <path>=<data>")
				continue
			}

			kvPart := parts[1]
			if !strings.Contains(kvPart, "=") {
				fmt.Println("Invalid format. Use 'set <path>=<data>'")
				continue
			}

			kvParts := strings.Split(kvPart, "=")
			if len(kvParts) != 2 {
				fmt.Println("Invalid format. Use 'set <path>=<data>'")
				continue
			}
			path, data := kvParts[0], kvParts[1]
			success, result := z.Create(data, path)
			fmt.Printf("Set %s: %v\n", map[bool]string{true: "successful", false: "failed"}[success], result)
		case "get":
			if len(parts) < 2 {
				fmt.Println("Usage: get <path>")
				continue
			}

			path := parts[1]
			success, result := z.Read(path)
			if success {
				fmt.Printf("Data at '%s': %v\n", path, result)
			} else {
				fmt.Printf("Error: %v\n", result)
			}
		case "ls":
			if len(parts) < 2 {
				path := z.basePath
				children, err := z.getChildren(path)
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					continue
				}
				fmt.Printf("Children of '%s': %v\n", path, children)
				continue
			}

			path := parts[1]
			children, err := z.getChildren(path)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			fmt.Printf("Children of '%s': %v\n", path, children)
		case "stat":
			_, result := z.Stats(nil)
			fmt.Printf("Stats: %v\n", result)
		case "delete":
			if len(parts) < 2 {
				fmt.Println("Usage: delete <path>")
				continue
			}

			path := parts[1]
			success, result := z.Delete(path)
			fmt.Printf("Delete %s: %v\n", map[bool]string{true: "successful", false: "failed"}[success], result)
		default:
			fmt.Printf("Unknown command: %s\n", action)
		}
	}
}

func (z *ZkCrud) getChildren(path string) ([]string, error) {
	if z.conn == nil {
		z.connect()
	}

	children, _, err := z.conn.Children(path)
	if err != nil {
		return nil, err
	}
	return children, nil
}
