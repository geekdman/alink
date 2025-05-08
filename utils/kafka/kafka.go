package kafka

//
//import (
//	"context"
//	"fmt"
//	"log"
//	"strconv"
//	"strings"
//
//	"alink/utils/base"
//	"github.com/segmentio/kafka-go"
//	"bufio"
//	"os"
//)
//
//type KafkaCrud struct {
//	ip     string
//	port   string
//	reader *kafka.Reader
//	writer *kafka.Writer
//}
//
//func NewKafkaCrud(info string) *KafkaCrud {
//	parts := strings.Split(info, ":")
//	if len(parts) != 2 {
//		log.Fatal("Invalid Kafka connection info. Use 'ip:port' format.")
//	}
//
//	return &KafkaCrud{
//		ip:   parts[0],
//		port: parts[1],
//	}
//}
//
//func (k *KafkaCrud) connect() {
//	address := fmt.Sprintf("%s:%s", k.ip, k.port)
//	k.reader = kafka.NewReader(kafka.ReaderConfig{
//		Brokers:  []string{address},
//		Topic:    "default_topic",
//		MinBytes: 10e3, // 10KB
//		MaxBytes: 10e6, // 10MB
//	})
//
//	k.writer = kafka.NewWriter(kafka.WriterConfig{
//		Brokers: []string{address},
//		Topic:   "default_topic",
//	})
//}
//
//func (k *KafkaCrud) Create(data map[string]interface{}, identifier string) (bool, interface{}) {
//	if k.writer == nil {
//		k.connect()
//	}
//
//	err := k.writer.WriteMessages(context.Background(), kafka.Message{
//		Value: []byte(fmt.Sprint(data)),
//	})
//	if err != nil {
//		return false, err
//	}
//
//	return true, "Message sent successfully"
//}
//
//func (k *KafkaCrud) Read(query interface{}) (bool, interface{}) {
//	if k.reader == nil {
//		k.connect()
//	}
//
//	count := 1
//	if query != nil {
//		count = query.(int)
//	}
//
//	messages := make([]string, 0)
//	for i := 0; i < count; i++ {
//		msg, err := k.reader.ReadMessage(context.Background())
//		if err != nil {
//			break
//		}
//		messages = append(messages, string(msg.Value))
//	}
//
//	return true, messages
//}
//
//func (k *KafkaCrud) Update(identifier string, newData map[string]interface{}) (bool, interface{}) {
//	return k.Create(newData, identifier)
//}
//
//func (k *KafkaCrud) Delete(identifier string) (bool, interface{}) {
//	return false, "Delete operation is not supported for Kafka"
//}
//
//func (k *KafkaCrud) Stats(query interface{}) (bool, interface{}) {
//	return true, "Statistics information is not directly available in the current Kafka client version."
//}
//
//func (k *KafkaCrud) Run() {
//	k.connect()
//	fmt.Printf("\n=== Kafka Shell (Connected to %s:%s) ===\n", k.ip, k.port)
//	fmt.Println("Commands: set, get, ls, stat, delete, help, exit")
//
//	scanner := bufio.NewScanner(os.Stdin)
//	for {
//		fmt.Print("kafka> ")
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
//			fmt.Println("  set <topic> <message>   - Send a message to a topic")
//			fmt.Println("  get <count>             - Consume messages from a topic")
//			fmt.Println("  ls                      - List available topics")
//			fmt.Println("  stat                    - Show Kafka cluster statistics")
//			fmt.Println("  delete <topic>          - Delete a topic")
//			fmt.Println("  exit                    - Exit the Kafka shell")
//			fmt.Println("  help                    - Show this help message")
//		case "set":
//			if len(parts) < 2 {
//				fmt.Println("Usage: set <topic> <message>")
//				continue
//			}
//
//			topicMsg := parts[1]
//			topic, message, found := strings.Cut(topicMsg, " ")
//			if !found {
//				fmt.Println("Invalid format. Use 'set <topic> <message>'")
//				continue
//			}
//
//			success, result := k.Create(map[string]interface{}{"message": message}, topic)
//			fmt.Printf("Set %s: %v\n", map[bool]string{true: "successful", false: "failed"}[success], result)
//		case "get":
//			if len(parts) < 2 {
//				fmt.Println("Usage: get <count>")
//				continue
//			}
//
//			count, err := strconv.Atoi(parts[1])
//			if err != nil {
//				fmt.Println("Invalid count. Please provide a valid integer.")
//				continue
//			}
//
//			success, result := k.Read(count)
//			if success {
//				fmt.Printf("Messages: %v\n", result)
//			} else {
//				fmt.Printf("Error: %v\n", result)
//			}
//		case "ls":
//			// 实现 ls 命令
//		case "stat":
//			success, result := k.Stats(nil)
//			fmt.Printf("Stats: %v\n", result)
//		case "delete":
//			if len(parts) < 2 {
//				fmt.Println("Usage: delete <topic>")
//				continue
//			}
//
//			topic := parts[1]
//			success, result := k.Delete(topic)
//			fmt.Printf("Delete %s: %v\n", map[bool]string{true: "successful", false: "failed"}[success], result)
//		default:
//			fmt.Printf("Unknown command: %s\n", action)
//		}
//	}
//}
