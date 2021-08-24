package service

import (
	"log"
	"log-agent/entity"
	"time"

	"gopkg.in/Shopify/sarama.v1"
)

type KafkaService interface {
	SendDataToKafka()
	SendDataToKafkaChan(topic, data string)
}

type kafkaService struct {
	Producer sarama.SyncProducer //生产者
	logChan chan *entity.KafkaLogMsg // 存放数据的管道，需要make使用
}



func NewKafkaService(addr []string, config *sarama.Config, maxChanSize int) KafkaService {
	c := setConfig(config)
	producer, err := sarama.NewSyncProducer(addr, c)
	if err != nil {
		log.Fatalln("err", err)
	}
	kafkaService := &kafkaService{
		Producer: producer,
		logChan: make(chan *entity.KafkaLogMsg, maxChanSize),
	}
	go kafkaService.SendDataToKafka()
	return kafkaService
}
func setConfig(c *sarama.Config) (config *sarama.Config) {
	c.Producer.RequiredAcks = sarama.WaitForAll
	c.Producer.Partitioner	= sarama.NewRandomPartitioner
	c.Producer.Return.Successes = true

	return c
}
// 把要发送的日志存到chan
func (s *kafkaService) SendDataToKafkaChan(topic, data string) {
	msg := &entity.KafkaLogMsg{
		Topic: topic,
		Data: data,
	}
	s.logChan <- msg

}

func (s *kafkaService) SendDataToKafka() {
	for {
		select {
			case msg :=<- s.logChan:
				logData := &sarama.ProducerMessage{}
				logData.Topic = msg.Topic
				logData.Value = sarama.StringEncoder(msg.Data)
				// 发布消息
				partition, offset, err := s.Producer.SendMessage(logData)
				if err != nil {
					log.Println("send fail, parition=", partition, ",offset = ", offset, "err=", err)
				}
				log.Println("发送成功：，partition=", partition, "offset=", offset)
			default:
				time.Sleep(time.Millisecond * 50)
		}

	}

	
}