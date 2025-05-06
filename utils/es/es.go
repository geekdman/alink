package es

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/your/repo/utils/base"
)

type EsCrud struct {
	ip      string
	port    string
	index   string
	es      *elasticsearch.Client
	baseURL string
}

func NewEsCrud(info string) *EsCrud {
	parts := strings.Split(info, ":")
	if len(parts) != 2 {
		log.Fatal("Invalid Elasticsearch connection info. Use 'ip:port' format.")
	}

	return &EsCrud{
		ip:      parts[0],
		port:    parts[1],
		baseURL: fmt.Sprintf("http://%s:%s", parts[0], parts[1]),
	}
}

func (e *EsCrud) connect() {
	var err error
	e.es, err = elasticsearch.NewClient(elasticsearch.Config{
		Addresses: []string{e.baseURL},
	})
	if err != nil {
		log.Fatalf("Error creating Elasticsearch client: %s", err)
	}
}

func (e *EsCrud) Create(data map[string]interface{}, identifier string) (bool, interface{}) {
	if e.es == nil {
		e.connect()
	}

	ctx := context.Background()
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(data); err != nil {
		return false, err
	}

	req := esapi.IndexRequest{
		Index:      e.index,
		DocumentID: identifier,
		Body:       strings.NewReader(buf.String()),
		Refresh:    "true",
	}

	res, err := req.Do(ctx, e.es)
	if err != nil {
		return false, err
	}
	defer res.Body.Close()

	if res.IsError() {
		return false, fmt.Sprintf("[%s] %s", res.Status(), res.String())
	}

	return true, "Document created successfully"
}

func (e *EsCrud) Read(query interface{}) (bool, interface{}) {
	// 实现读取操作
}

func (e *EsCrud) Update(identifier string, newData map[string]interface{}) (bool, interface{}) {
	// 实现更新操作
}

func (e *EsCrud) Delete(identifier string) (bool, interface{}) {
	// 实现删除操作
}

func (e *EsCrud) Stats(query interface{}) (bool, interface{}) {
	// 实现统计操作
}

func (e *EsCrud) Run() {
	e.connect()
	fmt.Printf("\n=== Elasticsearch Shell (Connected to %s:%s) ===\n", e.ip, e.port)
	fmt.Println("Commands: set, get, ls, stat, delete, help, exit")

	for {
		fmt.Print(fmt.Sprintf("es[%s]> ", e.index))
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
			fmt.Println("  set [id=]<key>=<value>   - Create or update a document")
			fmt.Println("  get <id>                 - Retrieve a document by ID")
			fmt.Println("  ls                       - List documents in the index")
			fmt.Println("  stat                     - Show cluster statistics")
			fmt.Println("  delete <id>              - Delete a document by ID")
			fmt.Println("  exit                     - Exit the Elasticsearch shell")
			fmt.Println("  help                     - Show this help message")
			continue
		}

		// 处理其他命令
	}
}