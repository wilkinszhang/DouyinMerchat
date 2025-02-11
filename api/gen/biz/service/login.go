package service

import (
	"DouyinMerchant/api/gen/biz/dal/mysql"
	user "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/user"
	"context"
	"errors"
	"gorm.io/gorm"
)

type LoginService struct {
	ctx context.Context
} // NewLoginService new LoginService
func NewLoginService(ctx context.Context) *LoginService {
	return &LoginService{ctx: ctx}
}

// Run create note info
func (s *LoginService) Run(req *user.LoginReq) (resp *user.LoginResp, err error) {
	// Finish your business logic.
	if err := validateLoginReq(req); err != nil {
		return nil, err
	}

	userModel, err := mysql.GetUserByEmail(s.ctx, mysql.DB, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("invalid email or password")
		}
		return nil, err
	}
	hashedInput := hashPassword(req.Password)
	if userModel.Passowrd != hashedInput {
		return nil, errors.New("invalid email or password")
	}
	return &user.LoginResp{
		UserId: int32(userModel.ID),
	}, nil
}

func validateLoginReq(req *user.LoginReq) error {
	if req.Email == "" || req.Password == "" {
		return errors.New("email and password are required")
	}
	return nil
}
