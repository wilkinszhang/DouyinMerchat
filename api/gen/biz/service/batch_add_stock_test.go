package service

import (
	stock "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/stock"
	"context"
	"testing"
)

func TestBatchAddStock_Run(t *testing.T) {
	ctx := context.Background()
	s := NewBatchAddStockService(ctx)
	// init req and assert value

	req := &stock.BatchAddStockReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
