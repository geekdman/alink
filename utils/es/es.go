package es

//
//import (
//	"context"
//	"bytes"
//	"encoding/json"
//	"fmt"
//	"log"
//	"strings"
//	"time"
//
//	"alink/utils/base"
//	"github.com/elastic/go-elasticsearch/v8"
//	"github.com/elastic/go-elasticsearch/v8/esapi"
//	"bufio"
//	"os"
//)
//
//type EsCrud struct {
//	ip      string
//	port    string
//	index   string
//	es      *elasticsearch.Client
//	baseURL string
//}
//
//func NewEsCrud(info string) *EsCrud {
//	parts := strings.Split(info, ":")
//	if len(parts) != 2 {
//		log.Fatal("Invalid Elasticsearch connection info. Use 'ip:port' format.")
//	}
//
//	return &EsCrud{
//		ip:      parts[0],
//		port:    parts[1],
//		baseURL: fmt.Sprintf("http://%s:%s", parts[0], parts[1]),
//		index:   "default_index",
//	}
//}
//
//func (e *EsCrud) connect() {
//	var err error
//	e.es, err = elasticsearch.NewClient(elasticsearch.Config{
//		Addresses: []string{e.baseURL},
//	})
//	if err != nil {
//		log.Fatalf("Error creating Elasticsearch client: %s", err)
//	}
//}
//
//func (e *EsCrud) Create(data map[string]interface{}, identifier string) (bool, interface{}) {
//	if e.es == nil {
//		e.connect()
//	}
//
//	ctx := context.Background()
//	var buf bytes.Buffer
//	if err := json.NewEncoder(&buf).Encode(data); err != nil {
//		return false, err
//	}
//
//	req := esapi.IndexRequest{
//		Index:      e.index,
//		DocumentID: identifier,
//		Body:       strings.NewReader(buf.String()),
//		Refresh:    "true",
//	}
//
//	res, err := req.Do(ctx, e.es)
//	if err != nil {
//		return false, err
//	}
//	defer res.Body.Close()
//
//	if res.IsError() {
//		return false, fmt.Sprintf("[%s] %s", res.Status(), res.String())
//	}
//
//	return true, "Document created successfully"
//}
//
//func (e *EsCrud) Read(query interface{}) (bool, interface{}) {
//	if e.es == nil {
//		e.connect()
//	}
//
//	ctx := context.Background()
//	var buf bytes.Buffer
//	if err := json.NewEncoder(&buf).Encode(query); err != nil {
//		return false, err
//	}
//
//	req := esapi.SearchRequest{
//		Index: []string{e.index},
//		Body:  &buf,
//	}
//
//	res, err := req.Do(ctx, e.es)
//	if err != nil {
//		return false, err
//	}
//	defer res.Body.Close()
//
//	if res.IsError() {
//		return false, fmt.Sprintf("[%s] %s", res.Status(), res.String())
//	}
//
//	var result map[string]interface{}
//	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
//		return false, err
//	}
//
//	return true, result
//}
//
//func (e *EsCrud) Update(identifier string, newData map[string]interface{}) (bool, interface{}) {
//	return e.Create(newData, identifier)
//}
//
//func (e *EsCrud) Delete(identifier string) (bool, interface{}) {
//	if e.es == nil {
//		e.connect()
//	}
//
//	ctx := context.Background()
//	req := esapi.DeleteRequest{
//		Index:      e.index,
//		DocumentID: identifier,
//		Refresh:    "true",
//	}
//
//	res, err := req.Do(ctx, e.es)
//	if err != nil {
//		return false, err
//	}
//	defer res.Body.Close()
//
//	if res.IsError() {
//		return false, fmt.Sprintf("[%s] %s", res.Status(), res.String())
//	}
//
//	return true, "Document deleted successfully"
//}
//
//func (e *EsCrud) Stats(query interface{}) (bool, interface{}) {
//	if e.es == nil {
//		e.connect()
//	}
//
//	ctx := context.Background()
//	req := esapi.ClusterStatsRequest{}
//
//	res, err := req.Do(ctx, e.es)
//	if err != nil {
//		return false, err
//	}
//	defer res.Body.Close()
//
//	if res.IsError() {
//		return false, fmt.Sprintf("[%s] %s", res.Status(), res.String())
//	}
//
//	var result map[string]interface{}
//	if err := json.NewDecoder(res.Body).Decode(&result); err != nil {
//		return false, err
//	}
//
//	return true, result
//}
//
//func (e *EsCrud) Run() {
//	e.connect()
//	fmt.Printf("\n=== Elasticsearch Shell (Connected to %s:%s) ===\n", e.ip, e.port)
//	fmt.Println("Commands: set, get, ls, stat, delete, help, exit")
//
//	scanner := bufio.NewScanner(os.Stdin)
//	for {
//		fmt.Print(fmt.Sprintf("es[%s]> ", e.index))
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
//			fmt.Println("  set [id=]<key>=<value>   - Create or update a document")
//			fmt.Println("  get <id>                 - Retrieve a document by ID")
//			fmt.Println("  ls                       - List documents in the index")
//			fmt.Println("  stat                     - Show cluster statistics")
//			fmt.Println("  delete <id>              - Delete a document by ID")
//			fmt.Println("  exit                     - Exit the Elasticsearch shell")
//			fmt.Println("  help                     - Show this help message")
//		case "set":
//			if len(parts) < 2 {
//				fmt.Println("Usage: set [id=]<key>=<value>")
//				fmt.Println("Example: set id=1 name=John age=30")
//				fmt.Println("Or: set name=John age=30")
//				continue
//			}
//
//			dataPart := parts[1]
//			docID := ""
//			data := make(map[string]interface{})
//
//			// 检查是否有指定的文档 ID
//			if strings.HasPrefix(dataPart, "id=") {
//				idPart, dataPart := strings.SplitN(dataPart, " ", 2)
//				docID = strings.Split(idPart, "=")[1]
//			}
//
//			// 解析键值对
//			kvPairs := strings.Split(dataPart, ",")
//			for _, pair := range kvPairs {
//				if strings.Contains(pair, "=") {
//					key, value := strings.Split(pair, "=", 2)
//					data[key] = value
//				}
//			}
//
//			success, result := e.Create(data, docID)
//			fmt.Printf("Set %s: %v\n", map[bool]string{true: "successful", false: "failed"}[success], result)
//		case "get":
//			if len(parts) < 2 {
//				fmt.Println("Usage: get <id>")
//				continue
//			}
//
//			docID := parts[1]
//			success, result := e.Read(docID)
//			if success {
//				fmt.Printf("Document '%s': %v\n", docID, result)
//			} else {
//				fmt.Printf("Error: %v\n", result)
//			}
//		case "ls":
//			// 实现 ls 命令
//		case "stat":
//			success, result := e.Stats(nil)
//			fmt.Printf("Stats: %v\n", result)
//		case "delete":
//			if len(parts) < 2 {
//				fmt.Println("Usage: delete <id>")
//				continue
//			}
//
//			docID := parts[1]
//			success, result := e.Delete(docID)
//			fmt.Printf("Delete %s: %v\n", map[bool]string{true: "successful", false: "failed"}[success], result)
//		default:
//			fmt.Printf("Unknown command: %s\n", action)
//		}
//	}
//}
