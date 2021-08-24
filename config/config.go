package config

import (
	"log"
	"time"
	"gopkg.in/ini.v1"
)
type Config interface {
	GetKafkaAddr() string
	GetKafkaTopic() string
	GetKafkaChanMaxsize() int
	GetTaillogPath() string
	GetEtcdAddr() string
	GetEtcdTaillogKey() string
	GetEtcdTimeout() time.Duration
}
type config struct {
	kafkaAddress string
	kafkaTopic string
	kafkaChanMaxsize int
	logFilepath string
	etcdAddress string
	etcdTimeOut time.Duration
	etcdTaillogKey string
}

func LoadConfig() Config {
	cfg, err := ini.Load("./config.ini")
	if err != nil {
		log.Fatalln("err", err)
	}
	return &config{
		kafkaAddress: cfg.Section("kafka").Key("address").String(),
		kafkaTopic: cfg.Section("kafka").Key("topic").String(),
		kafkaChanMaxsize: cfg.Section("kafka").Key("chann_size").MustInt(10000),
		logFilepath: cfg.Section("taillog").Key("filepath").String(),
		etcdAddress: cfg.Section("etcd").Key("address").String(),
		etcdTaillogKey : cfg.Section("etcd").Key("log_agent_key").String(),
		etcdTimeOut: cfg.Section("etcd").Key("timeout").MustDuration(5),
	}
}
func (this *config) GetKafkaAddr() string {
	return this.kafkaAddress
}

func (this *config) GetKafkaTopic() string {
	return this.kafkaTopic
}

func (this *config) GetKafkaChanMaxsize() int {
	return this.kafkaChanMaxsize
}

func (this *config) GetTaillogPath() string {
	return this.logFilepath
}
func (this *config) GetEtcdAddr() string {
	return this.etcdAddress
}
func (this *config) GetEtcdTaillogKey() string {
	return this.etcdTaillogKey
}
func (this *config) GetEtcdTimeout() time.Duration {
	return this.etcdTimeOut
}
