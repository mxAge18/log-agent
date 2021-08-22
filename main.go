package main

import (
	"log-agent/config"
	"log-agent/controller"
	"log-agent/service"

	"github.com/hpcloud/tail"
	"gopkg.in/Shopify/sarama.v1"
)
var(
	Config config.Config = config.LoadConfig()
	KafkaService service.KafkaService = service.NewKafkaService([]string{Config.GetKafkaAddr()}, 
		sarama.NewConfig())
	TailService service.LogService = service.NewLogService(Config.GetTaillogPath(), tail.Config{
		ReOpen: true,
		Follow: true,
		Location: &tail.SeekInfo{Offset: 0, Whence: 2},
		MustExist: false,
		Poll: true,
	})
	LogController controller.LogController = controller.NewLogController(KafkaService, TailService)

)
func run() {
	LogController.SendLogToKafka(Config.GetKafkaTopic())
}

func main() {
	// 加载配置

	// 运行服务
	run()
}