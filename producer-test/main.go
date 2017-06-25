package main

import (
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	pb "github.com/tkanos/kafka-consumer/myService"
)

func main() {

	//addresses of available kafka braokers
	brokers := []string{"localhost:9092"}

	producer, err := NewProducer(brokers)

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}

	defer producer.Close()

	partition, offset, err := SendMessage(producer, EncodeMessage("1", "Felipe"), "my-service")
	fmt.Printf("%v %v %v", partition, offset, err)

}

// New Producer ...
func NewProducer(brokers []string) (sarama.SyncProducer, error) {
	//setup relevant config info
	//config := sarama.NewConfig()
	//config.Producer...... to have more information about configuration se https://godoc.org/github.com/Shopify/sarama#Config
	producer, err := sarama.NewSyncProducer(brokers, nil) //NewAsyncProducer(brokers, config)

	return producer, err
}

// SendMessage ...
func SendMessage(producer sarama.SyncProducer, msg []byte, topic string) (int32, int64, error) {

	message := &sarama.ProducerMessage{
		Topic:     topic,
		Partition: -1,
		Value:     sarama.ByteEncoder(msg),
	}
	partition, offset, err := producer.SendMessage(message)

	return partition, offset, err
}

// EncodeMessage ...
func EncodeMessage(id, name string) []byte {
	msg := &pb.MyServiceMessage{
		ID:   id,
		Name: name,
	}

	data, _ := proto.Marshal(msg)

	return data

}
