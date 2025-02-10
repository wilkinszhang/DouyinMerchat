package mysql

import (
	"DouyinMerchant/api/gen/biz/model"
	"context"
	"gorm.io/gorm"
)

type User struct {
	*model.User
}

func CreateUser(ctx context.Context, db *gorm.DB, user *User) error {
	return db.WithContext(ctx).Create(user.User).Error
}

func GetUserByEmail(ctx context.Context, db *gorm.DB, email string) (*User, error) {
	var user User
	if err := db.WithContext(ctx).Where("email=?", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
