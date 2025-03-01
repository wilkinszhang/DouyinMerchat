package service

import (
	stock "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/stock"
	"context"
)

type GetStockService struct {
	ctx context.Context
} // NewGetStockService new GetStockService
func NewGetStockService(ctx context.Context) *GetStockService {
	return &GetStockService{ctx: ctx}
}

// Run create note info
func (s *GetStockService) Run(req *stock.GetStockReq) (resp *stock.StockResp, err error) {
	// Finish your business logic.

	return
}
