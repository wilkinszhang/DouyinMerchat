package service

import (
	stock "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/stock"
	"context"
	"testing"
)

func TestDecreaseStock_Run(t *testing.T) {
	ctx := context.Background()
	s := NewDecreaseStockService(ctx)
	// init req and assert value

	req := &stock.DecreaseStockReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
