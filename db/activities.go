package db

import "time"

// Activity represents a message record in the database
type Activity struct {
	ID              int       `gorm:"primaryKey"`
	ActivityID      string    `gorm:"type:char(36);not null"`
	MessageTemplate string    `gorm:"type:varchar(255);not null"`
	ScheduledTime   time.Time `gorm:"type:datetime;not null"`
	Filename        string    `gorm:"type:varchar(255);not null"`
	Status          int       `gorm:"type:int;not null"`
}

// CreateActivity a new message
func CreateActivity(activityID, messageTemplate string, scheduledTime time.Time, filename string) (*Activity, error) {
	msg := Activity{
		ActivityID:      activityID,
		MessageTemplate: messageTemplate,
		ScheduledTime:   scheduledTime,
		Filename:        filename,
	}
	result := db.Create(&msg)
	return &msg, result.Error
}

func QueryNewActivity(limit int) ([]*Activity, error) {
	var activities []*Activity
	result := db.Where("status = ?", 0).Limit(limit).Find(&activities)
	if result.Error != nil {
		return activities, nil
	}
	return activities, result.Error
}

func UpdateActivityStatus(activityIDs []int) error {
	result := db.Model(&Activity{}).
		Where("id IN ?", activityIDs).Update("status", 1)
	return result.Error
}
