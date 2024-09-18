# iMessage

消息活动系统

## 使用方式
* #### HTTP Server
```shell
go run main.go -s server
```
接收并存储用户csv，以及消息模版等信息，csv格式如下
```shell
姓名,号码
```
上传文件方式
```shell
curl -X POST 'http://localhost:8080/upload?scheduled_time=2024-09-10%2000:00:00&message_template=dsfdsafdsfsd123099888222' -F "file=@demo.csv"
```

* #### Message Server
```shell
go run main.go -s message
```
解析csv，并将用户、消息数据入库

* #### Producer Server
```shell
go run main.go -s producer
```
生产者服务，按配置好的消息发送时间，将消息投入队列

* #### Consumer Server
```shell
go run main.go -s consumer
```
消费者服务，消费队列中的任务，发送消息并修改消息状态

## 组建配置
```go
// constant/net.go

package constant

const (
	QueueType   = queue.KafkaQueueType
	QueueAddr   = "127.0.0.1:9092"
	QueuePasswd = ""
	QueueKey    = "message"
)

const (
	MysqlUser   = "root"
	MysqlHost   = "127.0.0.1"
	MysqlPasswd = "Qwer1234!"
	MysqlDB     = "message"
	MysqlPort   = 3306
)

```