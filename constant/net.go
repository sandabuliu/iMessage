package constant

import "iMessage/client/queue"

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
