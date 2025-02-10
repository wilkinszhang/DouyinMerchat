package service

import (
	auth "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/auth"
	"context"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"os"
)

type VerifyTokenByRPCService struct {
	ctx context.Context
} // NewVerifyTokenByRPCService new VerifyTokenByRPCService
func NewVerifyTokenByRPCService(ctx context.Context) *VerifyTokenByRPCService {
	return &VerifyTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *VerifyTokenByRPCService) Run(req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	// Finish your business logic.
	if req.Token == "" {
		return &auth.VerifyResp{Res: false}, errors.New("empty token")
	}

	token, err := jwt.Parse(req.Token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil {
		return &auth.VerifyResp{Res: false}, err
	}
	//把token claims转化成MapClaims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//检查user_id在claims map里面是否存在
		if _, ok := claims["user_id"]; !ok {
			return &auth.VerifyResp{Res: false}, errors.New("invalid token claims")
		}
		return &auth.VerifyResp{Res: true}, nil
	}

	return &auth.VerifyResp{Res: false}, errors.New("invalid token")

	return
}
