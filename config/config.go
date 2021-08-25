package config

import (
	"fmt"
	"log"
	"log-agent/utils"
	"time"

	"gopkg.in/ini.v1"
)
type Config interface {
	GetKafkaAddr() string
	GetKafkaTopic() string
	GetKafkaChanMaxsize() int
	GetEtcdAddr() string
	GetEtcdTaillogKey() string
	GetEtcdTimeout() time.Duration
}
type config struct {
	KafkaConf `ini:"kafka"`
	EtcdConf `ini:"etcd"`
}
type KafkaConf struct {
	Address string `ini:"address"`
	Topic string `ini:"topic"`
	ChanMaxSize int `ini:"chann_size"`
}

type EtcdConf struct {
	Address string `ini:"address"`
	TimeOut time.Duration `ini:"timeout"`
	TaillogKey string `ini:"log_agent_key"`
}

func LoadConfig(path string) Config {
	var cfg config
	err := ini.MapTo(&cfg, path)
	if err != nil {
		log.Fatalln("load config.ini fail, err=", err)
	}
	log.Println("config load success")
	ip := utils.GetLocalIp()
	cfg.EtcdConf.TaillogKey = fmt.Sprintf(cfg.EtcdConf.TaillogKey, ip)
	return &cfg
}

func (this *config) GetKafkaAddr() string {
	return this.KafkaConf.Address
}

func (this *config) GetKafkaTopic() string {
	return this.KafkaConf.Topic
}

func (this *config) GetKafkaChanMaxsize() int {
	return this.KafkaConf.ChanMaxSize
}

func (this *config) GetEtcdAddr() string {
	return this.EtcdConf.Address
}

func (this *config) GetEtcdTaillogKey() string {
	return this.EtcdConf.TaillogKey
}

func (this *config) GetEtcdTimeout() time.Duration {
	return this.EtcdConf.TimeOut
}
