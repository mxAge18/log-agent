package main

import (
	"log-agent/controller"
	"log-agent/service"

	"github.com/hpcloud/tail"
	"gopkg.in/Shopify/sarama.v1"
)
var(
	KafkaService service.KafkaService = service.NewKafkaService([]string{"127.0.0.1:9092"}, 
		sarama.NewConfig())
	TailService service.LogService = service.NewLogService("my.log", tail.Config{
		ReOpen: true,
		Follow: true,
		Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll: true,
	})
	LogController controller.LogController = controller.NewLogController(KafkaService, TailService)
)
func run() {
	LogController.SendLogToKafka("web_log")
}
func main() {
	run()
}