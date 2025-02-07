package model

import (
	"context"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	UserId   int    `gorm:"uniqueIndex;type:int not null"`
	Password string `gorm:"type:varchar(64) not null"`
}

func Create(ctx context.Context, db *gorm.DB, user *User) error {
	return db.WithContext(ctx).Create(user).Error
}

func GetByUserId(ctx context.Context, db *gorm.DB, userid int32) (*User, error) {
	var user User
	err := db.WithContext(ctx).Where("userid = ?", userid).First(&user).Error
	return &user, err
}
