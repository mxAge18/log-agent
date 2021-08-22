package controller

import (
	"log-agent/service"
	"time"
)

type LogController interface {
	SendLogToKafka(topic string)
}

type logController struct {

	kafkaService service.KafkaService
	logService service.LogService
}

func NewLogController(Kafka service.KafkaService, LogService service.LogService) LogController {
	return &logController{
		kafkaService: Kafka,
		logService: LogService,
	}
}

func (c *logController) SendLogToKafka(topic string) {
	for {
		select {
			case line := <-c.logService.ReadLine():
				c.kafkaService.SendDataToKafka(topic, line.Text)
			default:
				time.Sleep(time.Second)

		}
	}

}