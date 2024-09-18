package main

import (
	"context"
	"flag"
	"iMessage/client/queue"
	"iMessage/constant"
	"iMessage/db"
	"iMessage/lib/activity"
	"iMessage/routes"
)

func init() {
	err := queue.InitQueue(constant.QueueType, constant.QueueAddr, constant.QueuePasswd, constant.QueueKey)
	if err != nil {
		panic(err)
	}

	err = db.InitDB(constant.MysqlUser, constant.MysqlPasswd, constant.MysqlHost, constant.MysqlDB, constant.MysqlPort)
	if err != nil {
		panic(err)
	}
}

func main() {
	typeFlag := flag.String("s", "server", "server/message/producer/consumer")
	flag.Parse()

	var err error
	ctx := context.Background()
	if *typeFlag == "server" {
		r := routes.SetupRouter() // 加载自定义路由
		err = r.Run(":8080")
	} else if *typeFlag == "message" {
		err = activity.Start(ctx, activity.CsvToUsers)
	} else if *typeFlag == "producer" {
		err = activity.Start(ctx, activity.Producer)
	} else if *typeFlag == "consumer" {
		err = activity.Start(ctx, activity.Consumer)
	}

	if err != nil {
		panic(err)
	}
}
