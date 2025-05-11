// utils/tikv/tikv.go
package tikv

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/chzyer/readline"
	"github.com/tikv/client-go/v2/txnkv"
)

type TiKVCrud struct {
	ip   string
	port string
	kv   *txnkv.Client
}

func NewTikvCrud(info string) *TiKVCrud {
	client, err := txnkv.NewClient([]string{info})

	if err != nil {
		log.Fatal("Error creating TiKV connection:", err)
	}
	parts := strings.Split(info, ":")

	return &TiKVCrud{
		ip:   parts[0],
		port: parts[1],
		kv:   client,
	}
}

func (t *TiKVCrud) Create(data interface{}, identifier string) (bool, interface{}) {
	txn, err := t.kv.Begin()
	if err != nil {
		return false, err
	}

	key := []byte(identifier)
	value := []byte(fmt.Sprint(data))

	err = txn.Set(key, value)
	if err != nil {
		txn.Rollback()
		return false, err
	}

	err = txn.Commit(context.Background())
	if err != nil {
		return false, err
	}

	return true, fmt.Sprintf("Data inserted at %s", identifier)
}

func (t *TiKVCrud) Read(query interface{}) (bool, interface{}) {
	key := []byte(query.(string))

	txn, err := t.kv.Begin()
	if err != nil {
		return false, err
	}
	defer txn.Rollback()

	value, err := txn.Get(context.Background(), key)
	if err != nil {
		return false, err
	}

	if value == nil {
		return false, fmt.Errorf("key not found")
	}

	return true, string(value)
}

func (t *TiKVCrud) Update(identifier string, newData interface{}) (bool, interface{}) {
	return t.Create(newData, identifier)
}

func (t *TiKVCrud) Delete(identifier string) (bool, interface{}) {
	txn, err := t.kv.Begin()
	if err != nil {
		return false, err
	}

	key := []byte(identifier)

	err = txn.Delete(key)
	if err != nil {
		txn.Rollback()
		return false, err
	}

	err = txn.Commit(context.Background())
	if err != nil {
		return false, err
	}

	return true, fmt.Sprintf("Data deleted from %s", identifier)
}

func (t *TiKVCrud) Stats() (bool, interface{}) {
	//stats, err := t.kv.GetClusterInfo(context.Background())
	//if err != nil {
	//	return false, err
	//}
	//
	//统计信息 := make(map[string]interface{})
	//统计信息["ClusterID"] = stats.ClusterID
	//统计信息["StoreCount"] = stats.Stores
	//统计信息["RegionCount"] = stats.RegionCount
	//t.kv.
	return true, "xxx"
}

func (t *TiKVCrud) Close() {
	t.kv.Close()
}

func (t *TiKVCrud) Run() {
	rl, err := readline.NewEx(&readline.Config{
		Prompt:          fmt.Sprintf("tikv[%s:%s]> ", t.ip, t.port),
		HistoryFile:     "/tmp/tikv_readline.tmp",
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	if err != nil {
		fmt.Printf("Error initializing readline: %v\n", err)
		return
	}
	defer rl.Close()

	fmt.Printf("\n=== TiKV Shell (Connected to %s:%s) ===\n", t.ip, t.port)
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
			t.Close()
			fmt.Println("Returning to main shell...")
			return
		case "help":
			fmt.Println("\nAvailable Commands:")
			fmt.Println("  set <key>=<data> - Create or update a key-value pair")
			fmt.Println("  get <key>        - Retrieve data by key")
			fmt.Println("  ls [ startkey [ endkey]]             - List keys in a prefix (not supported in this version)")
			fmt.Println("  stat             - Show TiKV cluster statistics")
			fmt.Println("  delete <key>     - Delete a key")
			fmt.Println("  exit             - Exit the TiKV shell")
			fmt.Println("  help             - Show this help message")
		case "set":
			if len(parts) < 2 {
				fmt.Println("Usage: set <key>=<data>")
				continue
			}

			kvPart := parts[1]
			if !strings.Contains(kvPart, "=") {
				fmt.Println("Invalid format. Use 'set <key>=<data>'")
				continue
			}

			kvParts := strings.Split(kvPart, "=")
			if len(kvParts) != 2 {
				fmt.Println("Invalid format. Use 'set <key>=<data>'")
				continue
			}
			key, data := kvParts[0], kvParts[1]
			success, result := t.Create(data, key)
			fmt.Printf("Set %s: %v\n", map[bool]string{true: "successful", false: "failed"}[success], result)
		case "get":
			if len(parts) < 2 {
				fmt.Println("Usage: get <key>")
				continue
			}

			key := parts[1]
			success, result := t.Read(key)
			if success {
				fmt.Printf("Data from '%s': %v\n", key, result)
			} else {
				fmt.Printf("Error: %v\n", result)
			}
		case "ls":
			parts := strings.SplitN(cmd, " ", 3) // Split into up to 3 parts to handle start and end keys
			if len(parts) < 2 {
				fmt.Println("Usage: ls <startkey> [<endkey>]")
				continue
			}

			startKey := parts[1]
			var endKey *string

			if len(parts) == 3 {
				endKeyStr := parts[2]
				endKey = &endKeyStr
			}

			t.ListKeys(startKey, endKey)

		case "stat":
			success, result := t.Stats()
			if success {
				jsonData, _ := json.MarshalIndent(result, "", "  ")
				fmt.Printf("TiKV Stats:\n%s\n", jsonData)
			} else {
				fmt.Printf("Error getting TiKV stats: %v\n", result)
			}
		case "delete":
			if len(parts) < 2 {
				fmt.Println("Usage: delete <key>")
				continue
			}

			key := parts[1]
			success, result := t.Delete(key)
			fmt.Printf("Delete %s: %v\n", map[bool]string{true: "successful", false: "failed"}[success], result)
		default:
			fmt.Printf("Unknown command: %s\n", action)
		}
	}
}

// 列出以什么开头的key
func (t *TiKVCrud) ListKeys(startKey string, endKey *string) {
	txn, err := t.kv.Begin()
	if err != nil {
		fmt.Printf("Error beginning transaction: %v\n", err)
		return
	}
	defer txn.Rollback()

	// 如果 endKey 未指定，则假定为无限大（即列出从 startKey 开始的所有键）
	var end []byte
	if endKey != nil {
		end = []byte(*endKey)
	} else {
		end = []byte("\xff\xff\xff\xff\xff\xff\xff\xff")
	}

	// 转换 startKey 为空字符串时的处理
	var start []byte
	if startKey != "" {
		start = []byte(startKey)
	} else {
		start = nil // 使用 nil 表示从头开始
	}

	// 扫描键空间中从 start 到 end 之间的所有键
	iter, err := txn.Iter(start, end)
	if err != nil {
		fmt.Printf("Error creating iterator: %v\n", err)
		return
	}
	defer iter.Close()

	fmt.Printf("Keys in the specified range:\n")
	for iter.Valid() {
		key := iter.Key()
		fmt.Println(string(key))
		if err := iter.Next(); err != nil {
			fmt.Printf("Error iterating keys: %v\n", err)
			return
		}
	}
}