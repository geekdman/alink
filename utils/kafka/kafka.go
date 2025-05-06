package kafka

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/your/repo/utils/base"
	"github.com/segmentio/kafka-go"
)

type KafkaCrud struct {
	ip     string
	port   string
	reader *kafka.Reader
	writer *kafka.Writer
}

func NewKafkaCrud(info string) *KafkaCrud {
	parts := strings.Split(info, ":")
	if len(parts) != 2 {
		log.Fatal("Invalid Kafka connection info. Use 'ip:port' format.")
	}

	return &KafkaCrud{
		ip:   parts[0],
		port: parts[1],
	}
}

func (k *KafkaCrud) connect() {
	address := fmt.Sprintf("%s:%s", k.ip, k.port)
	k.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{address},
		Topic:    "default_topic",
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})

	k.writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{address},
		Topic:   "default_topic",
	})
}

func (k *KafkaCrud) Create(topic, message string) (bool, interface{}) {
	if k.writer == nil {
		k.connect()
	}

	err := k.writer.WriteMessages(context.Background(), kafka.Message{
		Value: []byte(message),
	})
	if err != nil {
		return false, err
	}

	return true, "Message sent successfully"
}

func (k *KafkaCrud) Read(count int) (bool, interface{}) {
	if k.reader == nil {
		k.connect()
	}

	messages := make([]string, 0)
	for i := 0; i < count; i++ {
		msg, err := k.reader.ReadMessage(context.Background())
		if err != nil {
			break
		}
		messages = append(messages, string(msg.Value))
	}

	return true, messages
}

func (k *KafkaCrud) Update(topic, message string) (bool, interface{}) {
	return k.Create(topic, message)
}

func (k *KafkaCrud) Delete(topic string) (bool, interface{}) {
	// 实现删除操作
}

func (k *KafkaCrud) Stats(query interface{}) (bool, interface{}) {
	// 实现统计操作
}

func (k *KafkaCrud) Run() {
	k.connect()
	fmt.Printf("\n=== Kafka Shell (Connected to %s:%s) ===\n", k.ip, k.port)
	fmt.Println("Commands: set, get, ls, stat, delete, help, exit")

	for {
		fmt.Print("kafka> ")
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
			fmt.Println("  set <topic> <message>   - Send a message to a topic")
			fmt.Println("  get <count>             - Consume messages from a topic")
			fmt.Println("  ls                      - List available topics")
			fmt.Println("  stat                    - Show Kafka cluster statistics")
			fmt.Println("  delete <topic>          - Delete a topic")
			fmt.Println("  exit                    - Exit the Kafka shell")
			fmt.Println("  help                    - Show this help message")
			continue
		}

		// 处理其他命令
	}
}