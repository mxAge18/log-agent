package service

import (
	"context"
	"log"
	"log-agent/entity"
	"time"

	"github.com/hpcloud/tail"
)

// 利用tail 加载日志

type LogService interface {
	ReadLine() (<-chan *tail.Line)
	GetFileName() string
	initTailFile()
	SendDataToKafka()
	CancelFunc() context.CancelFunc
	LogContext() context.Context
}

type logService struct {
	filename string
	path string
	topic string
	logTails *tail.Tail
	config tail.Config
	err error
	Ctx context.Context // 为了变更配置项后，关闭正在运行的向kafka发数据的goroutine
	LogCancelFunc context.CancelFunc
	kafkaService KafkaService
}

func NewLogService(TailConfig tail.Config, LogEntity *entity.TaillogConfig, kafkaService KafkaService) LogService {
	ctx, cancel := context.WithCancel(context.Background())
	this := &logService{
		config: TailConfig,
		path: LogEntity.Path,
		topic: LogEntity.Topic,
		kafkaService: kafkaService,
		Ctx: ctx,
		LogCancelFunc: cancel,
	}
	this.initTailFile()
	go this.SendDataToKafka()
	return this
}

func (this *logService) initTailFile() {
	this.logTails, this.err = tail.TailFile(this.path, this.config)
	if this.err != nil {
		log.Fatalln("open file fail", this.err)
	}
}
func (this *logService) ReadLine() (<-chan *tail.Line) {
	return this.logTails.Lines
	
}
func (this *logService) SendDataToKafka() {
	for {
		// log数据存到kafka的消息chan中
		select {
			case <- this.Ctx.Done():
				log.Println("task 退出,topic",this.topic, "path", this.path)
				return
			case line := <-this.logTails.Lines:
				this.kafkaService.SendDataToKafkaChan(this.topic, line.Text)
				log.Println("send data to chan,topic=", this.topic, ",text=", line.Text)
			default:
				time.Sleep(time.Second)
		}
	}
}
func (this *logService) GetFileName() string {
	return this.filename
}

func (this *logService) CancelFunc() context.CancelFunc {
	return this.LogCancelFunc
}

func (this *logService) LogContext() context.Context {
	return this.Ctx
}
