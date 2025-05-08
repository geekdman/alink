package zk

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"log"
	"strings"
	"time"
)

type ZkCrud struct {
	ip     string
	port   string
	conn   *zk.Conn
	basePath string
}

func NewZkCrud(info string) *ZkCrud {
	parts := strings.Split(info, ":")
	if len(parts) != 2 {
		log.Fatal("Invalid ZooKeeper connection info. Use 'ip:port' format.")
	}

	return &ZkCrud{
		ip: parts[0],
		port: parts[1],
		basePath: "/",
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

func (z *ZkCrud) Create(data map[string]interface{}, identifier string) (bool, interface{}) {
	if z.conn == nil {
		z.connect()
	}

	path := z.basePath + identifier
	acl := zk.WorldACL(zk.PermAll)

	_, err := z.conn.Create(path, []byte(fmt.Sprint(data)), 0, acl)
	if err != nil {
		return false, err
	}

	return true, fmt.Sprintf("Node created at %s", path)
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

func (z *ZkCrud) Update(identifier string, newData map[string]interface{}) (bool, interface{}) {
	return z.Create(newData, identifier)
}

func (z *ZkCrud) Delete(identifier string) (bool, interface{}) {
	if z.conn == nil {
		z.connect()
	}

	path := z.basePath + identifier
	err := z.conn.Delete(path, -1)
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
			fmt.Println("  set <path>=<data> - Create or update a znode")
			fmt.Println("  get <path>         - Retrieve data from a znode")
			fmt.Println("  ls [path]          - List children of a znode")
			fmt.Println("  stat               - Show ZooKeeper server statistics")
			fmt.Println("  delete <path>      - Delete a znode")
			fmt.Println("  exit               - Exit the ZooKeeper shell")
			fmt.Println("  help               - Show this help message")
			continue
		}

		if action == "set" {
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
			success, result := z.Create(map[string]interface{}{"data": data}, path)
			fmt.Printf("Set %s: %v\n", map[bool]string{true: "successful", false: "failed"}[success], result)
		} else if action == "get" {
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
		} else if action == "ls" {
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
		} else if action == "stat" {
			_, result := z.Stats(nil)
			fmt.Printf("Stats: %v\n", result)
		} else if action == "delete" {
			if len(parts) < 2 {
				fmt.Println("Usage: delete <path>")
				continue
			}

			path := parts[1]
			success, result := z.Delete(path)
			fmt.Printf("Delete %s: %v\n", map[bool]string{true: "successful", false: "failed"}[success], result)
		} else {
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