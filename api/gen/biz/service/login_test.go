package service

import (
	user "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/user"
	"context"
	"testing"
)

func TestLogin_Run(t *testing.T) {
	ctx := context.Background()
	s := NewLoginService(ctx)
	// init req and assert value

	req := &user.LoginReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
	testValue := "123"
	t.Logf("helloworld: %v", testValue)
}
