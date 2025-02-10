package model

import (
	"context"
	"gorm.io/gorm"
	"time"
)

type User struct {
	ID        uint   `gorm:"primarykey"`
	Email     string `gorm:"uniqueIndex;size:255;not null"`
	Passowrd  string `gorm:"size:64;not null"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func Create(ctx context.Context, db *gorm.DB, user *User) error {
	return db.WithContext(ctx).Create(user).Error
}

func GetByUserId(ctx context.Context, db *gorm.DB, userid int32) (*User, error) {
	var user User
	err := db.WithContext(ctx).Where("userid = ?", userid).First(&user).Error
	return &user, err
}
