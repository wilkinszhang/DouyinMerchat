package service

import (
	"DouyinMerchant/api/gen/biz/dal/mysql"
	"DouyinMerchant/api/gen/biz/dal/redis"
	"DouyinMerchant/api/gen/biz/model"
	stock "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/stock"
	"context"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
)

type BatchAddStockService struct {
	ctx context.Context
}

// NewBatchAddStockService new BatchAddStockService
func NewBatchAddStockService(ctx context.Context) *BatchAddStockService {
	return &BatchAddStockService{ctx: ctx}
}

// Run create note info
func (s *BatchAddStockService) Run(req *stock.BatchAddStockReq) (resp *stock.StockResp, err error) {
	// Initialize response
	resp = &stock.StockResp{
		Success: true,
		Stocks:  make(map[uint32]*stock.StockInfo),
	}

	// Validate request
	if len(req.Items) == 0 {
		resp.Success = false
		resp.Message = "No items provided"
		return resp, nil
	}

	// Update Redis cache first
	pipe := redis.RedisClient.Pipeline()

	// Track stock changes for database update
	stockChanges := make(map[uint32]int32)

	for _, item := range req.Items {
		// Validate item
		if item.Num <= 0 {
			klog.Warnf("Skipping invalid stock quantity for product %d: %d", item.ProductId, item.Num)
			continue
		}

		stockKey := fmt.Sprintf("prod_{%d}_stock", item.ProductId)

		// Increment stock in Redis
		pipe.IncrBy(s.ctx, stockKey, int64(item.Num))

		// Track for database update
		stockChanges[item.ProductId] = item.Num
	}

	// Execute Redis operations
	if _, err := pipe.Exec(s.ctx); err != nil {
		klog.Errorf("Failed to update Redis stock: %v", err)
		resp.Success = false
		resp.Message = "Failed to update stock cache"
		return resp, err
	}

	// Update database asynchronously
	go func() {
		// Create context for the goroutine
		ctx := context.Background()

		err := model.BatchAddStock(mysql.DB, ctx, stockChanges)
		if err != nil {
			klog.Errorf("Failed to update database stock: %v", err)

			// Attempt to rollback Redis changes
			rollbackPipe := redis.RedisClient.Pipeline()
			for productID, quantity := range stockChanges {
				rollbackPipe.DecrBy(ctx, fmt.Sprintf("prod_{%d}_stock", productID), int64(quantity))
			}

			if _, rollbackErr := rollbackPipe.Exec(ctx); rollbackErr != nil {
				klog.Errorf("Failed to rollback Redis stock changes: %v", rollbackErr)
			}
		}
	}()

	// Get updated stock for response
	stockService := NewGetStockService(s.ctx)
	productIDs := make([]uint32, 0, len(req.Items))

	for _, item := range req.Items {
		productIDs = append(productIDs, item.ProductId)
	}

	stockResp, err := stockService.Run(&stock.GetStockReq{ProductIds: productIDs})
	if err != nil {
		klog.Warnf("Failed to get updated stock: %v", err)
		// Continue with empty stocks map
	} else {
		resp.Stocks = stockResp.Stocks
	}

	return resp, nil
}
