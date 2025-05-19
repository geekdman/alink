// utils/tikv/tikv.go
package tikv

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/chzyer/readline"
	mylog "github.com/pingcap/log"
	"github.com/tikv/client-go/v2/txnkv"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

type TiKVCrud struct {
	ip   string
	port string
	kv   *txnkv.Client
}

type Clusterinfo struct {
	pdstatus  interface{}
	tikvstatus interface{}
	clusterinfo interface{}
}

func NewTikvCrud(endpoint string) *TiKVCrud {

	mylog.SetLevel(zapcore.ErrorLevel)
	client, err := txnkv.NewClient([]string{endpoint})

	if err != nil {
		log.Fatal("Error creating TiKV connection:", err)
	}
	parts := strings.Split(endpoint, ":")

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
	pdstatus, err1 := t.GetPDstatus()
	tikvstatus, err2 := t.GetTikvstatus()
	clusterinfo, err3 := t.GetClustertatus()

	if err2 != nil || err3 != nil {
		return false, fmt.Errorf("error retrieving cluster info: pdstatus error: %v, tikvstatus error: %v, clusterinfo error: %v", err1, err2, err3)
	}
	// 创建一个 map 存储集群信息
	clusterMap := make(map[string]interface{})
	clusterMap["pdstatus"] = pdstatus
	clusterMap["tikvstatus"] = tikvstatus
	clusterMap["clusterinfo"] = clusterinfo

	return true, clusterMap
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
			fmt.Println("  ls <startkey>    - List keys in a prefix (not supported in this version)")
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
			parts := strings.SplitN(cmd, " ", 2) // Split into up to 2 parts
			if len(parts) < 2 {
				fmt.Println("Usage: ls <startkey>")
				continue
			}

			startKey := parts[1]
			t.ListKeys(startKey)

		case "stat":
			success, result := t.Stats()
			if success {
				jsonData, err := json.MarshalIndent(result, "", "  ")
				if err != nil {
					fmt.Printf("Error formatting TiKV stats: %v\n", err)
				} else {
					fmt.Printf("TiKV Cluster Status:\n%s\n", jsonData)
				}
			} else {
				fmt.Printf("Error getting TiKV cluster status: %v\n", result)
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

// 列出从 startKey 开始的所有键，直到遇到比 startKey 字典序大的键为止
func (t *TiKVCrud) ListKeys(startKey string) {
	txn, err := t.kv.Begin()
	if err != nil {
		fmt.Printf("Error beginning transaction: %v\n", err)
		return
	}
	defer txn.Rollback()

	start := []byte(startKey)

	// 创建迭代器，从 startKey 开始，到 startKey 的下一个字典序键结束
	end := []byte(startKey)
	// 如果 startKey 是空字符串，从头开始
	if len(startKey) == 0 {
		start = nil
		end = []byte("")
	} else {
		// 找到 startKey 的下一个字典序键作为 end
		for i := len(end) - 1; i >= 0; i-- {
			if end[i] < 0xff {
				end[i]++
				end = end[:i+1]
				break
			}
		}
		if len(end) == 0 {
			end = []byte{0x00} // 如果 startKey 全是 0xff，则 end 设置为 0x00
		}
	}

	iter, err := txn.Iter(start, end)
	if err != nil {
		fmt.Printf("Error creating iterator: %v\n", err)
		return
	}
	defer iter.Close()

	// 检查是否没有键
	if !iter.Valid() {
		fmt.Printf("No keys found starting with '%s'.\n", startKey)
		return
	}

	// 打印所有从 startKey 开始的键
	fmt.Println(strings.Repeat("==", 10))
	for iter.Valid() {
		key := iter.Key()
		fmt.Println(string(key))
		if err := iter.Next(); err != nil {
			fmt.Printf("Error iterating keys: %v\n", err)
			return
		}
	}
}

// 访问pd 的api
func (t *TiKVCrud) Request(apiurl string) (interface{}, error) {
	resp, err := http.Get(apiurl)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PD: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get cluster status: HTTP status %d", resp.StatusCode)
	}

	byteData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var res map[string]interface{}
	if err := json.Unmarshal(byteData, &res); err != nil {
		return nil, fmt.Errorf("failed to parse JSON response: %v", err)
	}

	return &res, nil
}

type PDStatus struct {
	Name       string   `json:"name"`
	MemberId   float64  `json:"member_id"`
	ClientUrls []string `json:"client_urls"`
	Health     bool     `json:"health"`
}

type TikvStatus struct {
	Count  int `json:"count"`
	Stores []struct {
		Store struct {
			Id             int    `json:"id"`
			Address        string `json:"address"`
			Version        string `json:"version"`
			PeerAddress    string `json:"peer_address"`
			StatusAddress  string `json:"status_address"`
			GitHash        string `json:"git_hash"`
			StartTimestamp int    `json:"start_timestamp"`
			DeployPath     string `json:"deploy_path"`
			LastHeartbeat  int64  `json:"last_heartbeat"`
			NodeState      int    `json:"node_state"`
			StateName      string `json:"state_name"`
		} `json:"store"`
		Status struct {
			Capacity     string  `json:"capacity"`
			Available    string  `json:"available"`
			UsedSize     string  `json:"used_size"`
			LeaderCount  int     `json:"leader_count"`
			LeaderWeight int     `json:"leader_weight"`
			LeaderScore  int     `json:"leader_score"`
			LeaderSize   int     `json:"leader_size"`
			RegionCount  int     `json:"region_count"`
			RegionWeight int     `json:"region_weight"`
			RegionScore  float64 `json:"region_score"`
			RegionSize   int     `json:"region_size"`
			SlowScore    int     `json:"slow_score"`
			SlowTrend    struct {
				CauseValue  int `json:"cause_value"`
				CauseRate   int `json:"cause_rate"`
				ResultValue int `json:"result_value"`
				ResultRate  int `json:"result_rate"`
			} `json:"slow_trend"`
			StartTs         time.Time `json:"start_ts"`
			LastHeartbeatTs time.Time `json:"last_heartbeat_ts"`
			Uptime          string    `json:"uptime"`
		} `json:"status"`
	} `json:"stores"`
}

type CluserStatus struct {
	Id           int64 `json:"id"`
	MaxPeerCount int   `json:"max_peer_count"`
}
// 获取pd 状态
// http://%s:%s/pd/api/v1/health
func (t *TiKVCrud) GetPDstatus() (interface{}, error) {
	pdInfo, err := t.Request(fmt.Sprintf("http://%s:%s/pd/api/v1/health", t.ip, t.port))
	if err != nil {
		return nil, err
	}

	// 检查是否为数组
	pdArray, ok := pdInfo.([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format from PD API: not an array")
	}

	// 遍历数组查找包含 name 和 health 的元素
	var pdStatus []map[string]interface{}
	for _, item := range pdArray {
		itemMap, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if name, ok := itemMap["name"].(string); ok {
			if health, ok := itemMap["health"].(bool); ok {
				pdStatus = append(pdStatus, map[string]interface{}{
					"name":    name,
					"health":  health,
				})
				continue
			}
		}
	}

	if pdStatus == nil {
		return nil, fmt.Errorf("PD status information not found in the response")
	}
	fmt.Println("%v",pdStatus)
	return pdStatus, nil
}

// 获取tikv 状态
// http://%s:%s/pd/api/v1/stores
func (t *TiKVCrud) GetTikvstatus() (interface{}, error) {
	clusterTikvInfo, err := t.Request(fmt.Sprintf("http://%s:%s/pd/api/v1/stores", t.ip, t.port))
	if err != nil {
		return nil, err
	}

	// 将结果转换为 map
	tikvMap, ok := clusterTikvInfo.(*map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected response format from TiKV API")
	}

	// 提取 stores 信息
	stores, ok := (*tikvMap)["stores"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("unexpected format for stores in TiKV API response")
	}

	// 提取每个 TiKV 实例的 address、version 和 state_name
	var tikvStores []map[string]interface{}
	for _, store := range stores {
		storeMap, ok := store.(map[string]interface{})
		if !ok {
			continue
		}

		// 提取 store 信息
		storeInfo, ok := storeMap["store"].(map[string]interface{})
		if !ok {
			continue
		}

		address, ok := storeInfo["address"].(string)
		if !ok {
			address = "N/A"
		}

		version, ok := storeInfo["version"].(string)
		if !ok {
			version = "N/A"
		}

		stateName, ok := storeInfo["state_name"].(string)
		if !ok {
			stateName = "N/A"
		}

		tikvStores = append(tikvStores, map[string]interface{}{
			"address":    address,
			"version":    version,
			"state_name": stateName,
		})
	}

	return tikvStores, nil
}

// 获取 集群信息
// http://%s:%s/pd/api/v1/cluster
func (t *TiKVCrud) GetClustertatus() (interface{}, error) {
	clusterPdInfo, err := t.Request(fmt.Sprintf("http://%s:%s/pd/api/v1/cluster", t.ip, t.port))
	if err != nil {
		return nil, err
	}
	return clusterPdInfo, nil
}