package service

import (
	"log"

	"gopkg.in/Shopify/sarama.v1"
)

type KafkaService interface {
	SendDataToKafka(topic string, data string) (err error)
}

type kafkaService struct {
	Producer sarama.SyncProducer //生产者
}

func setConfig(c *sarama.Config) (config *sarama.Config) {
	c.Producer.RequiredAcks = sarama.WaitForAll
	c.Producer.Partitioner	= sarama.NewRandomPartitioner
	c.Producer.Return.Successes = true
	return c
}
func NewKafkaService(addr []string, config *sarama.Config) KafkaService {
	c := setConfig(config)
	producer, err := sarama.NewSyncProducer(addr, c)
	if err != nil {
		log.Fatalln("err", err)
	}
	return &kafkaService{
		Producer: producer,
	}
}

func (service *kafkaService) SendDataToKafka(topic string, data string) (err error) {
	msg := &sarama.ProducerMessage{}
	msg.Topic = topic
	msg.Value = sarama.StringEncoder(data)

	// 发布消息
	partition, offset, err := service.Producer.SendMessage(msg)
	if err != nil {
		log.Println("send success, parition=", partition, ",offset = ", offset)
		return err
	}
	log.Println("发送成功：，partition=", partition, "offset=", offset)
	return nil
	
}