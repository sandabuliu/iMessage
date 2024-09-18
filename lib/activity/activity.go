package activity

import (
	"context"
	"fmt"
	"iMessage/client/message"
	"iMessage/client/queue"
	"iMessage/db"
	"iMessage/utils"
	"log"
	"strconv"
	"time"
)

func formatMessages(activity *db.Activity, users []*db.User) map[int]string {
	result := map[int]string{}
	for _, user := range users {
		// todo: 补充模版生成业务逻辑
		result[user.ID] = activity.MessageTemplate
	}
	return result
}

func CsvToUsers(ctx context.Context, limit int) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}()

	activities, err := db.QueryNewActivity(limit)
	if err != nil {
		return 0, err
	}
	log.Printf("query activity num: %d\n", len(activities))

	counter := 0
	for _, act := range activities {
		users := map[string]string{}
		err = utils.CsvLoop(act.Filename, func(cols []string) error {
			if len(cols) != 2 {
				return fmt.Errorf("invalid column count: %d", len(cols))
			}
			users[cols[1]] = cols[0]
			return nil
		})
		if err != nil {
			return counter, err
		}
		log.Printf("user from csv, user num: %d\n", len(users))

		_, err := db.InsertUsers(users)
		if err != nil {
			return counter, err
		}
		userPhones := make([]string, 0)
		for phone, _ := range users {
			userPhones = append(userPhones, phone)
		}
		dbUsers, err := db.QueryUsersByPhones(userPhones)
		if err != nil {
			return counter, err
		}
		log.Printf("insert users, num: %d\n", len(dbUsers))

		err = db.InsertMessages(act.ID, act.ScheduledTime, formatMessages(act, dbUsers))
		if err != nil {
			return counter, err
		}
		log.Printf("insert messages, num: %d\n", len(dbUsers))

		err = db.UpdateActivityStatus([]int{act.ID})
		if err != nil {
			return counter, err
		}
		log.Printf("update activity status\n")
		counter++
	}
	return counter, nil
}

func Producer(ctx context.Context, limit int) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Println("Recovered from panic:", r)
		}
	}()

	messages, err := db.QuerySendMessages(limit)
	if err != nil {
		return 0, err
	}
	log.Printf("query new messages, num: %d\n", len(messages))

	counter := 0
	for _, msg := range messages {
		err = queue.Client.Produce(ctx, fmt.Sprintf("%d", msg.ID))
		if err != nil {
			return counter, err
		}
		err = db.UpdateMessagesStatus([]int{msg.ID}, db.StatusQueue)
		if err != nil {
			return counter, err
		}
	}

	return counter, nil
}

func Consumer(ctx context.Context, limit int) (int, error) {
	defer func() {
		if r := recover(); r != nil {
			time.Sleep(10 * time.Second)
			log.Println("Recovered from panic:", r)
		}
	}()

	messageIds := make([]int, 0)
	for i := 0; i < limit; i++ {
		log.Printf("consume start\n")
		msgId, err := queue.Client.Consume(ctx)
		if err != nil {
			log.Printf("consume failed, error: %+v, num: %d", err, i)
			break
		}
		log.Printf("consume message_id: %d\n", msgId)
		messageId, _ := strconv.Atoi(msgId)
		messageIds = append(messageIds, messageId)
	}
	if len(messageIds) == 0 {
		return 0, nil
	}

	messages, err := db.QueryMessagesByIds(messageIds)
	if err != nil {
		log.Printf("query messages by id failed, error: %+v, msgIds: %+v", err, messageIds)
		return 0, err
	}
	err = message.Send(ctx, messages)
	if err != nil {
		log.Printf("send message failed, error: %+v, msgIds: %+v", err, messageIds)
		return 0, err
	}

	err = db.UpdateMessagesStatus(messageIds, db.StatusSent)
	if err != nil {
		log.Printf("update status failed, error: %+v, msgIds: %+v", err, messageIds)
		return 0, err
	}
	return len(messages), err
}

func Start(ctx context.Context, process func(context.Context, int) (int, error)) error {
	limit := 10
	log.Println("start process")
	for {
		num, err := process(ctx, limit)
		if err != nil {
			log.Fatalf("process: %+v", err.Error())
		}

		if num < limit || err != nil {
			time.Sleep(10 * time.Second)
		}
	}
	return nil
}
