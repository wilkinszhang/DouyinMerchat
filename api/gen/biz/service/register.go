package service

import (
	"DouyinMerchant/api/gen/biz/dal/mysql"
	"DouyinMerchant/api/gen/biz/model"
	user "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/user"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
)

//{
//"email": "33@bilibili.com",
//"password": "33good",
//"confirm_password": "33good"
//}

type RegisterService struct {
	ctx context.Context
} // NewRegisterService new RegisterService
func NewRegisterService(ctx context.Context) *RegisterService {
	return &RegisterService{ctx: ctx}
}

// Run create note info
func (s *RegisterService) Run(req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	// Finish your business logic.
	//验证请求
	if err := validateRegisterReq(req); err != nil {
		return nil, err
	}
	//验证邮箱是否存在
	if _, err := mysql.GetUserByEmail(s.ctx, mysql.DB, req.Email); err == nil {
		return nil, errors.New("email already exists")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	hashedPassword := hashPassword(req.Password)

	userModel := &mysql.User{
		User: &model.User{
			Email:    req.Email,
			Passowrd: hashedPassword,
		},
	}

	//创建用户
	if err := mysql.CreateUser(s.ctx, mysql.DB, userModel); err != nil {
		return nil, err
	}
	return &user.RegisterResp{
		UserId: int32(userModel.ID),
	}, nil
}

func hashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return hex.EncodeToString(hash.Sum(nil))
}

func validateRegisterReq(req *user.RegisterReq) error {
	if req.Email == "" {
		return errors.New("email is required")
	}
	if req.Password == "" {
		return errors.New("password is required")
	}
	if req.Password != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}
	return nil
}
