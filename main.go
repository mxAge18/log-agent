package main

import (
	"fmt"
	"log-agent/config"
	"log-agent/controller"
	"log-agent/server_manager"
	"log-agent/service"
	"time"

	"github.com/hpcloud/tail"
	"gopkg.in/Shopify/sarama.v1"
)
var(
	Config config.Config = config.LoadConfig("./config/config.ini")
	KafkaService service.KafkaService = service.NewKafkaService(
		[]string{Config.GetKafkaAddr()}, 
		sarama.NewConfig(), 
		Config.GetKafkaChanMaxsize())

	TaillogManger *server_manager.TaillogManager = server_manager.NewTaillogManager(
		tail.Config{
			ReOpen: true,
			Follow: true,
			Location: &tail.SeekInfo{Offset: 0, Whence: 2},
			MustExist: false,
			Poll: true,},
			KafkaService,
		)

	EtcdService service.EtcdService = service.NewEtcdService(
		Config.GetEtcdAddr(),
		Config.GetEtcdTimeout() * time.Second)

	LogController controller.LogController = controller.NewLogController(
			EtcdService, 
			KafkaService, 
			TaillogManger, Config.GetEtcdTaillogKey())

)
func run() {
	LogController.SendMultiLogToKafka()
}

func main() {
	// 运行服务
	run()
}