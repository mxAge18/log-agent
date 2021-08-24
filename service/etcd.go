package service

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"log-agent/entity"
	"time"

	"go.etcd.io/etcd/client/v3"
)

type EtcdService interface {
	setConfig()
	initClient() (err error)
	ReadConfig(key string) (logConfigEntity []*entity.TaillogConfig, err error) // 从etcd中读取config
	PutConfig(key string, val string) (err error)
	WatchKey(key string, newConfigChan chan[]*entity.TaillogConfig)
}

type etcdService struct {

	address string
	connetTimeout time.Duration
	cli *clientv3.Client
	cfg clientv3.Config
}

func NewEtcdService(etcdAddress string, etcdTimeoutSet time.Duration) EtcdService {
	e := &etcdService{
		address: etcdAddress,
		connetTimeout: etcdTimeoutSet,
	}
	e.setConfig()
	e.initClient()
	return e
}

func (this *etcdService) setConfig() {
	this.cfg.Endpoints = []string{this.address}
	this.cfg.DialTimeout = this.connetTimeout
}
func (this *etcdService) initClient() (err error){
	this.cli, err = clientv3.New(this.cfg)
	if err != nil {
		// handle error!
		log.Fatalln("连接etcd错误：err", err)
		return 
	}
	return nil
}
func (this *etcdService) ReadConfig(key string) (logConfigEntity []*entity.TaillogConfig, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), this.connetTimeout * time.Second)
	resp, err := this.cli.Get(ctx, key)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	for _, ev := range resp.Kvs {
		
		json.Unmarshal(ev.Value, &logConfigEntity)
		if err != nil {
			log.Fatal(err)
			return 
		}
	}
	fmt.Println("获取到配置信息", resp.Kvs)
	return
}
func (this *etcdService) PutConfig(key string, val string) (err error) {
	ctx, cancel := context.WithTimeout(context.Background(), this.connetTimeout * time.Second)
	presp, err := this.cli.Put(ctx, key, val)
	cancel()
	if err != nil {
		log.Fatal(err)
	}
	
	fmt.Println(presp)
	return
}
func (this *etcdService) WatchKey(key string, newConfigChan chan[]*entity.TaillogConfig) {
	rch := this.cli.Watch(context.Background(), key)
	// 从监听的通道中尝试读取监听的信息
	for wresp := range rch {
		for _, ev := range wresp.Events {
			// 判断操作的类型，不是删除则解析ev.Kv.Value 然后存到管道中
			fmt.Println("检测到etcd存入的配置变化：,ev.value", ev.Kv.Value)
			var newConfig []*entity.TaillogConfig
			if ev.Type != clientv3.EventTypeDelete {
				err := json.Unmarshal(ev.Kv.Value, &newConfig)
				if err != nil {
					log.Fatalln("json unmarshal fail, watchkey func,", err)
					continue
				}
				newConfigChan <- newConfig
			}
		}
	}
}
