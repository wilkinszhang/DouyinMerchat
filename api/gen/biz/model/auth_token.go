package model

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"time"
)

type Auth_Token struct {
	gorm.Model
	UserId     int32  `gorm:"uniqueIndex;type:int32 not null"`
	Token      string `gorm:"type:varchar(512) not null"`
	CreateTime time.Time
	ExpireTime time.Time
	Status     int
}

func CreateToken(ctx context.Context, db *gorm.DB, auth_token *Auth_Token) error {
	return db.WithContext(ctx).Create(auth_token).Error
}

func CreateCacheToken(ctx context.Context, cacheClient *redis.Client, auth_token *Auth_Token) error {
	cacheKey := fmt.Sprintf("%s_%d", "token", auth_token.UserId)
	cacheClient.Set(ctx, cacheKey, 1, time.Hour)
	return nil
}
