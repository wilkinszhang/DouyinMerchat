package service

import (
	stock "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/stock"
	"context"
)

type PreReduceStockService struct {
	ctx context.Context
} // NewPreReduceStockService new PreReduceStockService
func NewPreReduceStockService(ctx context.Context) *PreReduceStockService {
	return &PreReduceStockService{ctx: ctx}
}

// Run create note info
func (s *PreReduceStockService) Run(req *stock.PreReduceStockReq) (resp *stock.StockResp, err error) {
	// Finish your business logic.

	return
}
