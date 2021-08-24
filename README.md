#log agent项目
##  环境搭建

- 本地测试开发，使用本人的go-docker-kits 快速启动kafka etcd zookeeper

## 业务逻辑

- go run main.go 运行
- 1 初始化etcd kafka连接
- 2 从etcd中读取特定key的value(logfile的配置信息，包含要使用什么topic发送到kafka，以及log文件的路径)
- 3 加载多个文件并开启，监听文件变化，读取文件并传到管道中
- 4 监听管道的内容，读取管道内的信息，发送到kafka中
- 5 建一个监听etcd配置的管道，变化时更新要加载的文件，关闭多余的goroutine,开启新的goroutine

## 遇到的问题

- 如何关闭goroutine
- 初始化要加载的内容

## 提高

- raft协议
  - 如何选举
  - 脑裂
  - 日志复制
- etcd watch实现方式，如何实现给client发送通知？
- context学习

## 下一步

- test 测试用例
- 统一代码风格
