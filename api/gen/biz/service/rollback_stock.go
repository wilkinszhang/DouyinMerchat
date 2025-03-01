package service

import (
	stock "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/stock"
	"context"
)

type RollbackStockService struct {
	ctx context.Context
} // NewRollbackStockService new RollbackStockService
func NewRollbackStockService(ctx context.Context) *RollbackStockService {
	return &RollbackStockService{ctx: ctx}
}

// Run create note info
func (s *RollbackStockService) Run(req *stock.RollbackStockReq) (resp *stock.StockResp, err error) {
	// Finish your business logic.

	return
}
