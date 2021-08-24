package server_manager

import (
	"fmt"
	"log-agent/entity"
	"log-agent/service"
	"time"

	"github.com/hpcloud/tail"
)

type TaillogManager struct {
	Taillog map[string]service.LogService
	newConfigChan chan []*entity.TaillogConfig
	Config tail.Config
	kafka service.KafkaService
}

func NewTaillogManager(config tail.Config, kafkaService service.KafkaService) *TaillogManager {
	tail := &TaillogManager{
		Config: config,
		newConfigChan: make(chan []*entity.TaillogConfig), // 无缓冲区的通道，，没有值阻塞
		Taillog: make(map[string]service.LogService, 16),
		kafka: kafkaService,
	}
	go tail.Run()
	return tail
}

// 监听newconf chan
// 配置更新：新增 删除 修改等
func (t *TaillogManager) Run() {
	for {
		select {
		case newConf := <- t.newConfigChan:
			// 新增
			fmt.Println("newConf")
			var keysMap map[string]int = make(map[string]int, 32)
			for _, taillogNewConf := range newConf {
				key := fmt.Sprintf("%s_%s", taillogNewConf.Topic, taillogNewConf.Path)
				_, ok := t.Taillog[key]
				keysMap[key] = 1
				if ok {
					// 配置没变，不调整
					fmt.Println("检测到配置：key,", key)
					continue
				} else {
					// 不存在 新增
					fmt.Println("检测到配置新增：key,", key)
					t.Taillog[key] = service.NewLogService(t.Config, taillogNewConf, t.kafka)
				}
				
			}

			for i := range t.Taillog {
				_, ok := keysMap[i]
				if ok {
					continue
				} else {
					t.Taillog[i].CancelFunc()()
					delete(t.Taillog, i)
					fmt.Println("删除key,", i)
				}
			}

		default:
			time.Sleep(time.Millisecond * 50)
		}
	
	
	}
}

// 定义一个函数向newConfigChan发送数据
func (t *TaillogManager) GetNewConfChan() chan []*entity.TaillogConfig {
	return t.newConfigChan
}