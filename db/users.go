package db

import (
	"gorm.io/gorm/clause"
)

type User struct {
	ID    int    `gorm:"primaryKey"`
	Name  string `gorm:"size:100;not null"`
	Phone string `gorm:"size:15;not null"`
}

func QueryUsersByPhones(phones []string) ([]*User, error) {
	var users []*User
	result := db.Where("phone IN ?", phones).Find(&users)
	return users, result.Error
}

func QueryUsersByIds(ids []int) ([]*User, error) {
	var users []*User
	result := db.Where("id IN ?", ids).Find(&users)
	return users, result.Error
}

func InsertUsers(user map[string]string) ([]*User, error) {
	// 批量插入用户数据
	users := make([]*User, 0)
	for phone, name := range user {
		users = append(users, &User{Name: name, Phone: phone})
	}

	result := db.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "phone"}},
	}).Create(&users)
	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}
