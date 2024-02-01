package kafka

import (
	"alink/config"
	"fmt"
	kafka "github.com/segmentio/kafka-go"
)

type kafkaConn struct {
	Kafkaconn *kafka.Conn
}

// 获取kafka 连接
func GetKafkaConn()  *kafkaConn{
	kafkaconn := new(kafkaConn)
	hosts := config.Cfg.GetKafkaConfig()

	for i := range hosts {
		conn ,err := kafka.Dial("tcp",hosts[i])
		if err != nil {
			panic(err)
		}
		kafkaconn.Kafkaconn = conn
		break
	}
	//fmt.Println(conn.Server())
	return kafkaconn
}

func (c *kafkaConn)GetTopics()  {
	partitions, err := c.Kafkaconn.ReadPartitions()

	if err != nil {
		panic(err.Error())
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	for k := range m {
		fmt.Println(k)
	}

}