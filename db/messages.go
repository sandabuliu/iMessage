package db

import (
	"time"
)

const (
	StatusNew   = 0
	StatusQueue = 1
	StatusSent  = 2
)

type Message struct {
	ID         int    `gorm:"primaryKey"`
	Message    string `gorm:"type:text;not null"`
	UserID     int    `gorm:"not null"`
	ActivityID int    `gorm:"not null"`
	Status     int    `gorm:"type:int;not null"`
}

func InsertMessages(activityId int, sendTime time.Time, messageContents map[int]string) error {
	messages := make([]Message, 0)
	for userId, messageContent := range messageContents {
		message := Message{
			Message:    messageContent,
			UserID:     userId,
			ActivityID: activityId,
		}
		messages = append(messages, message)
	}
	if err := db.Create(&messages).Error; err != nil {
		return err
	}
	return nil
}

// QuerySendMessages 根据用户ID批量查询消息的函数
func QuerySendMessages(limit int) ([]*Message, error) {
	var messages []*Message
	now := time.Now()
	err := db.Joins("JOIN activities ON activities.id = messages.activity_id").
		Where("messages.status = ? AND activities.scheduled_time < ?", 0, now).
		Find(&messages).Error
	return messages, err
}

func QueryMessagesByIds(ids []int) ([]*Message, error) {
	var messages []*Message
	result := db.Where("id IN ?", ids).Find(&messages)
	return messages, result.Error
}

func UpdateMessagesStatus(messageIds []int, status int) error {
	result := db.Model(&Message{}).
		Where("id IN ? AND status < ?", messageIds, status).Update("status", status)
	return result.Error
}
