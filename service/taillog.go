package service

import (
	"log"

	"github.com/hpcloud/tail"
)

// 利用tail 加载日志

type LogService interface {
	ReadLine() (<-chan *tail.Line)
	GetFileName() string
}

type logService struct {
	Filename string
	LogTails *tail.Tail
}

func NewLogService(fileName string, config tail.Config) LogService {
	tails, err := tail.TailFile(fileName, config)
	if err != nil {
		log.Fatalln("err,", err)
		panic("logService fail")
	}

	return &logService {
		LogTails: tails,
		Filename: tails.Filename,
	}
}

func (this *logService) ReadLine() (<-chan *tail.Line) {
	return this.LogTails.Lines
}

func (this *logService) GetFileName() string {
	return this.Filename
}
