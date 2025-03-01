package service

import (
	stock "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/stock"
	"context"
)

type BatchAddStockService struct {
	ctx context.Context
} // NewBatchAddStockService new BatchAddStockService
func NewBatchAddStockService(ctx context.Context) *BatchAddStockService {
	return &BatchAddStockService{ctx: ctx}
}

// Run create note info
func (s *BatchAddStockService) Run(req *stock.BatchAddStockReq) (resp *stock.StockResp, err error) {
	// Finish your business logic.

	return
}
