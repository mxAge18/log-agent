package controller

import (
	"fmt"
	"log"
	"log-agent/entity"
	"log-agent/server_manager"
	"log-agent/service"
	"sync"
	"time"
)

type LogController interface {
	getConfigFromEtcd()
	// SendLogToKafkaChan(topic string, s service.LogService)
	SendMultiLogToKafka()
	PutDataToEtcd()
}

type logController struct {
	etcdService service.EtcdService
	kafkaService service.KafkaService
	taillogManger *server_manager.TaillogManager
	taillogConfig []*entity.TaillogConfig // 存放etcd读取到的log路径和topic
	tailogEtcdKey string
	err error
}

func NewLogController(Etcd service.EtcdService, Kafka service.KafkaService, TaillogManger *server_manager.TaillogManager, etcdLogKey string) LogController {
	control := logController{
		etcdService: Etcd,
		kafkaService: Kafka,
		taillogManger: TaillogManger,
		tailogEtcdKey: etcdLogKey,
	}
	// 从ETCD 拉取kafka topic的路径及topic信息 存到control中
	control.getConfigFromEtcd()

	return &control
}
func (c *logController) getConfigFromEtcd() {
	c.taillogConfig, c.err = c.etcdService.ReadConfig(c.tailogEtcdKey)
	if c.err != nil {
		log.Fatalln("read config from etcd fail", c.err)
	}
}
func (c *logController) SendMultiLogToKafka() {

	for _, logEntity := range c.taillogConfig {
		key := fmt.Sprintf("%s_%s", logEntity.Topic, logEntity.Path)
		c.taillogManger.Taillog[key] = service.NewLogService(c.taillogManger.Config, logEntity, c.kafkaService)
		fmt.Println("logPath is", logEntity.Path)
	}
	newConfigChan := c.taillogManger.GetNewConfChan()
	// 开启一个后台gorutine watch etcd的配置变化
	var wg sync.WaitGroup
	wg.Add(1)
	go c.etcdService.WatchKey(c.tailogEtcdKey, newConfigChan)
	wg.Wait()
	fmt.Println(c.taillogManger.Taillog)

}

func (c *logController) SendLogToKafkaChan(topic string, s service.LogService) {
	for {
		// log数据存到kafka的消息chan中
		select {
			case <- s.LogContext().Done():
				log.Println("task 退出,topic",topic, "filename", s.GetFileName())
				return
			case line := <-s.ReadLine():
				c.kafkaService.SendDataToKafkaChan(topic, line.Text)
				fmt.Println("send data to chan,topic=", topic, ",text=", line.Text)
			default:
				time.Sleep(time.Second)
		}
	}
}

func (c *logController)PutDataToEtcd() {
	// value := `[{"path":"./test2.log","topic":"test2_log"},{"path":"./test3.log","topic":"test3_log"},{"path":"./test1.log","topic":"test1_log"}]`
	value := `[{"path":"./test1.log","topic":"test1_log"}]`
	c.etcdService.PutConfig(c.tailogEtcdKey, value)
}

