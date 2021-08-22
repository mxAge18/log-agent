package config

import (
	"log"
	"gopkg.in/ini.v1"
)
type Config interface {
	GetKafkaAddr() string
	GetKafkaTopic() string
	GetTaillogPath() string
}
type config struct {
	kafkaAddress string
	kafkaTopic string
	logFilepath string
}

func LoadConfig() Config {
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		log.Fatalln("err", err)
	}
	return &config{
		kafkaAddress: cfg.Section("kafka").Key("address").String(),
		kafkaTopic: cfg.Section("kafka").Key("topic").String(),
		logFilepath: cfg.Section("taillog").Key("filepath").String(),
	}
}
func (this *config) GetKafkaAddr() string {
	return this.kafkaAddress
}

func (this *config) GetKafkaTopic() string {
	return this.kafkaTopic
}

func (this *config) GetTaillogPath() string {
	return this.logFilepath
}
