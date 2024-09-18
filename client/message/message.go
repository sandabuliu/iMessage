package message

import (
	"context"
	"encoding/json"
	"fmt"
	"iMessage/db"
)

func Send(ctx context.Context, messages []*db.Message) error {
	userIds := make([]int, 0)
	for _, message := range messages {
		userIds = append(userIds, message.UserID)
	}
	users, err := db.QueryUsersByIds(userIds)
	if err != nil {
		return err
	}

	userMap := make(map[int]*db.User)
	for _, user := range users {
		userMap[user.ID] = user
	}
	for _, message := range messages {
		msg, _ := json.Marshal(message)
		user, _ := json.Marshal(userMap[message.UserID])
		fmt.Printf("send: %s, user: %s\n", string(msg), string(user))
	}
	return nil
}
