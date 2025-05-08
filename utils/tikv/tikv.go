package tikv

//
//import (
//	"context"
//	"fmt"
//	"log"
//	"strings"
//
//	"alink/utils/base"
//	"github.com/tikv/client-go/v2"
//	"bufio"
//	"os"
//)
//
//type TikvCrud struct {
//	ip     string
//	port   string
//	client *tikv.Client
//}
//
//func NewTikvCrud(info string) *TikvCrud {
//	parts := strings.Split(info, ":")
//	if len(parts) != 2 {
//		log.Fatal("Invalid TiKV connection info. Use 'ip:port' format.")
//	}
//
//	return &TikvCrud{
//		ip:   parts[0],
//		port: parts[1],
//	}
//}
//
//func (t *TikvCrud) connect() {
//	pdAddrs := []string{fmt.Sprintf("%s:%s", t.ip, t.port)}
//	t.client = tikv.NewClient(pdAddrs, nil)
//}
//
//func (t *TikvCrud) Create(data map[string]interface{}, identifier string) (bool, interface{}) {
//	if t.client == nil {
//		t.connect()
//	}
//
//	txn := t.client.Begin()
//	if err := txn.Set([]byte(identifier), []byte(fmt.Sprint(data))); err != nil {
//		return false, err
//	}
//
//	if err := txn.Commit(context.Background()); err != nil {
//		return false, err
//	}
//
//	return true, "Key created successfully"
//}
//
//func (t *TikvCrud) Read(query interface{}) (bool, interface{}) {
//	if t.client == nil {
//		t.connect()
//	}
//
//	snapshot := t.client.Snapshot(tikv.NewVersion())
//	value, err := snapshot.Get(context.Background(), []byte(query.(string)))
//	if err != nil {
//		return false, err
//	}
//
//	if value == nil {
//		return false, "Key not found"
//	}
//
//	return true, string(value)
//}
//
//func (t *TikvCrud) Update(identifier string, newData map[string]interface{}) (bool, interface{}) {
//	return t.Create(newData, identifier)
//}
//
//func (t *TikvCrud) Delete(identifier string) (bool, interface{}) {
//	if t.client == nil {
//		t.connect()
//	}
//
//	txn := t.client.Begin()
//	if err := txn.Delete([]byte(identifier)); err != nil {
//		return false, err
//	}
//
//	if err := txn.Commit(context.Background()); err != nil {
//		return false, err
//	}
//
//	return true, "Key deleted successfully"
//}
//
//func (t *TikvCrud) Stats(query interface{}) (bool, interface{}) {
//	return true, "Statistics information is not directly available in the current TiKV client version."
//}
//
//func (t *TikvCrud) Run() {
//	t.connect()
//	fmt.Printf("\n=== TiKV Shell (Connected to %s:%s) ===\n", t.ip, t.port)
//	fmt.Println("Commands: set, get, ls, stat, delete, help, exit")
//
//	scanner := bufio.NewScanner(os.Stdin)
//	for {
//		fmt.Print("tikv> ")
//
//		// 使用 Scanner 读取整行输入
//		if !scanner.Scan() {
//			fmt.Println("Error reading input:", scanner.Err())
//			continue
//		}
//
//		cmd := scanner.Text()
//
//		cmd = strings.TrimSpace(cmd)
//		if cmd == "" {
//			continue
//		}
//
//		parts := strings.SplitN(cmd, " ", 2)
//		if len(parts) < 1 {
//			continue
//		}
//
//		action := parts[0]
//		switch action {
//		case "exit":
//			fmt.Println("Returning to main shell...")
//			return
//		case "help":
//			fmt.Println("\nAvailable Commands:")
//			fmt.Println("  set <key>=<value>        - Create or update a key-value pair")
//			fmt.Println("  get <key>                - Retrieve the value of a key")
//			fmt.Println("  ls [startkey [endkey]]   - List keys starting with startkey or within the range")
//			fmt.Println("  stat                     - Show statistics of the TiKV cluster")
//			fmt.Println("  delete <key>             - Delete a key-value pair")
//			fmt.Println("  exit                     - Exit the TiKV shell")
//			fmt.Println("  help                     - Show this help message")
//		case "set":
//			if len(parts) < 2 {
//				fmt.Println("Usage: set <key>=<value>")
//				continue
//			}
//
//			kvPart := parts[1]
//			if !strings.Contains(kvPart, "=") {
//				fmt.Println("Invalid format. Use 'set <key>=<value>'")
//				continue
//			}
//
//			key, value := strings.Split(kvPart, "=", 2)
//			success, result := t.Create(map[string]interface{}{"value": value}, key)
//			fmt.Printf("Set %s: %v\n", map[bool]string{true: "successful", false: "failed"}[success], result)
//		case "get":
//			if len(parts) < 2 {
//				fmt.Println("Usage: get <key>")
//				continue
//			}
//
//			key := parts[1]
//			success, result := t.Read(key)
//			if success {
//				fmt.Printf("Value of '%s': %v\n", key, result)
//			} else {
//				fmt.Printf("Error: %v\n", result)
//			}
//		case "ls":
//			// 实现 ls 命令
//		case "stat":
//			success, result := t.Stats(nil)
//			fmt.Printf("Stats: %v\n", result)
//		case "delete":
//			if len(parts) < 2 {
//				fmt.Println("Usage: delete <key>")
//				continue
//			}
//
//			key := parts[1]
//			success, result := t.Delete(key)
//			fmt.Printf("Delete %s: %v\n", map[bool]string{true: "successful", false: "failed"}[success], result)
//		default:
//			fmt.Printf("Unknown command: %s\n", action)
//		}
//	}
//}
