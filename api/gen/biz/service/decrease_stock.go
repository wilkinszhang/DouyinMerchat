package service

import (
	stock "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/stock"
	"context"
)

type DecreaseStockService struct {
	ctx context.Context
} // NewDecreaseStockService new DecreaseStockService
func NewDecreaseStockService(ctx context.Context) *DecreaseStockService {
	return &DecreaseStockService{ctx: ctx}
}

// Run create note info
func (s *DecreaseStockService) Run(req *stock.DecreaseStockReq) (resp *stock.StockResp, err error) {
	// Finish your business logic.

	return
}
