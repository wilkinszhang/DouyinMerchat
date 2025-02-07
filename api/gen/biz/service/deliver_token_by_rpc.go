package service

import (
	"DouyinMerchant/api/gen/biz/dal/mysql"
	"DouyinMerchant/api/gen/biz/dal/redis"
	"DouyinMerchant/api/gen/biz/model"
	auth "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/auth"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
	"time"
)

type DeliverTokenByRPCService struct {
	ctx context.Context
} // NewDeliverTokenByRPCService new DeliverTokenByRPCService
func NewDeliverTokenByRPCService(ctx context.Context) *DeliverTokenByRPCService {
	return &DeliverTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *DeliverTokenByRPCService) Run(req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	// Finish your business logic.
	//传入用户id，返回token
	if req.UserId <= 0 {
		return nil, errors.New("invalid user id")
	}
	claims := jwt.MapClaims{
		"user_id": req.UserId,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
		"iat":     time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := []byte(os.Getenv("JWT_SECRET"))

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return nil, err
	}
	return &auth.DeliveryResp{
		Token: tokenString,
	}, nil
}

func (s *DeliverTokenByRPCService) validateUser(userId int32) error {
	_, err := model.GetByUserId(s.ctx, mysql.DB, userId)
	if err != nil {
		return errors.New("user not found")
	}
	return nil
}

func (s *DeliverTokenByRPCService) storeToken(userid int32, tokenString string) error {
	return model.CreateToken(s.ctx, mysql.DB, &model.Auth_Token{
		UserId:     userid,
		Token:      tokenString,
		CreateTime: time.Now(),
		ExpireTime: time.Now().Add(time.Hour * 24 * 7),
	})
}

func (s *DeliverTokenByRPCService) cacheToken(userid int32, tokenString string) error {
	return model.CreateCacheToken(s.ctx, redis.RedisClient, &model.Auth_Token{
		UserId:     userid,
		Token:      tokenString,
		CreateTime: time.Now(),
		ExpireTime: time.Now().Add(time.Hour * 24 * 7),
	})
}
