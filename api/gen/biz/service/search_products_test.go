package service

import (
	product "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/product"
	"context"
	"testing"
)

func TestSearchProducts_Run(t *testing.T) {
	ctx := context.Background()
	s := NewSearchProductsService(ctx)
	// init req and assert value

	req := &product.SearchProductsReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test
	//TODO写单测和压测
}
