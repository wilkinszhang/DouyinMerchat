package service

import (
	"DouyinMerchant/api/gen/biz/dal/mysql"
	"DouyinMerchant/api/gen/biz/dal/redis"
	"DouyinMerchant/api/gen/biz/model"
	stock "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/stock"
	"context"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"strconv"
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
	// Initialize response
	resp = &stock.StockResp{
		Success: true,
		Stocks:  make(map[uint32]*stock.StockInfo),
	}

	// Check if rollback has already been processed
	rollbackKey := fmt.Sprintf("rollback:%s", req.OrderId)
	exists, err := redis.RedisClient.Exists(s.ctx, rollbackKey).Result()
	if err != nil {
		klog.Warnf("Failed to check if rollback is processed: %v", err)
		// Continue processing as if rollback hasn't been processed
	} else if exists > 0 {
		// Rollback already processed, return current stock levels
		stockService := NewGetStockService(s.ctx)
		productIDs := make([]uint32, 0, len(req.Items))

		for _, item := range req.Items {
			productIDs = append(productIDs, item.ProductId)
		}

		return stockService.Run(&stock.GetStockReq{ProductIds: productIDs})
	}

	// Process rollback in Redis first
	client := redis.RedisClient
	pipe := client.Pipeline()

	stockMap := make(map[uint32]int32)

	for _, item := range req.Items {
		productId := item.ProductId
		quantity := item.Num

		stockKey := fmt.Sprintf("prod_{%d}_stock", productId)
		frozenKey := fmt.Sprintf("prod_{%d}_frozen", productId)

		pipe.IncrBy(s.ctx, stockKey, int64(quantity))
		pipe.DecrBy(s.ctx, frozenKey, int64(quantity))

		stockMap[productId] = quantity
	}

	// Mark rollback as processed (persist indefinitely)
	pipe.Set(s.ctx, rollbackKey, 1, 0)

	if _, err := pipe.Exec(s.ctx); err != nil {
		klog.Errorf("Failed to update Redis for rollback: %v", err)
		resp.Success = false
		resp.Message = "Failed to rollback stock in cache"
		return resp, nil
	}

	// Process rollback in database asynchronously
	go func() {
		// Create context for the goroutine
		ctx := context.Background()

		tx := mysql.DB.Begin()
		success := true

		for productId, quantity := range stockMap {
			err := model.RollbackStock(tx, ctx, productId, quantity)
			if err != nil {
				klog.Errorf("Failed to rollback stock in database for product %d: %v", productId, err)
				success = false
				break
			}
		}

		if !success {
			tx.Rollback()
			// No need to rollback Redis as it's already updated
			return
		}

		if err := tx.Commit().Error; err != nil {
			klog.Errorf("Failed to commit rollback transaction: %v", err)
		}
	}()

	for productId := range stockMap {
		stockKey := fmt.Sprintf("prod_{%d}_stock", productId)
		frozenKey := fmt.Sprintf("prod_{%d}_frozen", productId)

		stockValue, err := client.Get(s.ctx, stockKey).Result()
		frozenValue, err2 := client.Get(s.ctx, frozenKey).Result()

		if err == nil && err2 == nil {
			stock, _ := strconv.ParseInt(stockValue, 10, 32)
			frozen, _ := strconv.ParseInt(frozenValue, 10, 32)

			resp.Stocks[productId] = &stock.StockInfo{
				ProductId:   productId,
				Stock:       int32(stock),
				FrozenStock: int32(frozen),
			}
		}
	}

	return resp, nil
	return
}
