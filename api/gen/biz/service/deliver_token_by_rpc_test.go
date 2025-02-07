package service

import (
	auth "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/auth"
	"context"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestDeliverTokenByRPC_Run(t *testing.T) {
	ctx := context.Background()
	s := NewDeliverTokenByRPCService(ctx)
	// init req and assert value

	req := &auth.DeliverTokenReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
	tests := []struct {
		name    string
		userID  int32
		wantErr bool
	}{
		{
			name:    "valid user id",
			userID:  123,
			wantErr: false,
		},
		{
			name:    "invalid user id",
			userID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &auth.DeliverTokenReq{
				UserId: tt.userID,
			}

			resp, err := s.Run(req)

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, resp.Token)

				// 可选:验证token
				token, err := jwt.Parse(resp.Token, func(token *jwt.Token) (interface{}, error) {
					return []byte(os.Getenv("JWT_SECRET")), nil
				})
				assert.NoError(t, err)
				assert.True(t, token.Valid)

				claims := token.Claims.(jwt.MapClaims)
				assert.Equal(t, float64(tt.userID), claims["user_id"])
			}
		})
	}
}
