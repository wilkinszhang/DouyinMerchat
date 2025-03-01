package service

import (
	"DouyinMerchant/api/gen/biz/dal/mysql"
	"DouyinMerchant/api/gen/biz/dal/redis"
	"DouyinMerchant/api/gen/biz/model"
	"DouyinMerchant/api/gen/kitex_gen/douyin_merchant/stock"
	"fmt"
	"github.com/cloudwego/kitex/pkg/klog"
	"strconv"
	"time"
)

type DecreaseStockService struct {
	ctx context.Context
}

// NewDecreaseStockService new DecreaseStockService
func NewDecreaseStockService(ctx context.Context) *DecreaseStockService {
	return &DecreaseStockService{ctx: ctx}
}

// Run create note info
func (s *DecreaseStockService) Run(req *stock.DecreaseStockReq) (resp *stock.StockResp, err error) {
	// Initialize response
	resp = &stock.StockResp{
		Success: true,
		Stocks:  make(map[uint32]*stock.StockInfo),
	}

	// Check if order has already been processed
	orderKey := fmt.Sprintf("order_processed:%s", req.OrderId)
	exists, err := redis.RedisClient.Exists(s.ctx, orderKey).Result()
	if err != nil {
		klog.Warnf("Failed to check if order is processed: %v", err)
		// Continue processing as if order hasn't been processed
	} else if exists > 0 {
		// Order already processed, return current stock levels
		return getStockFromCache(s.ctx, req.Items), nil
	}

	// Process order items
	client := redis.RedisClient
	stockMap := make(map[uint32]int32)

	for _, item := range req.Items {
		stockMap[item.ProductId] = item.Num
	}

	// Process stock reduction in database
	tx := mysql.DB.Begin()
	success := true

	for productId, quantity := range stockMap {
		err := model.ReduceStock(tx, s.ctx, productId, quantity)
		if err != nil {
			klog.Errorf("Failed to decrease stock for product %d: %v", productId, err)
			success = false
			break
		}
	}

	if !success {
		tx.Rollback()
		resp.Success = false
		resp.Message = "Failed to decrease stock"
		return resp, nil
	}

	if err := tx.Commit().Error; err != nil {
		klog.Errorf("Failed to commit stock transaction: %v", err)
		resp.Success = false
		resp.Message = "Failed to commit stock changes"
		return resp, nil
	}

	// Update Redis cache
	pipe := client.Pipeline()

	for productId, quantity := range stockMap {
		frozenKey := fmt.Sprintf("prod_{%d}_frozen", productId)
		pipe.DecrBy(s.ctx, frozenKey, int64(quantity))
	}

	// Mark order as processed (with 24-hour expiry)
	pipe.Set(s.ctx, orderKey, 1, 24*time.Hour)

	if _, err := pipe.Exec(s.ctx); err != nil {
		klog.Errorf("Failed to update Redis: %v", err)
		// Continue as database is the source of truth
	}

	// Get updated stock for response
	return getStockFromCache(s.ctx, req.Items), nil
}

func getStockFromCache(ctx context.Context, items []*stock.StockItem) *stock.StockResp {
	resp := &stock.StockResp{
		Success: true,
		Stocks:  make(map[uint32]*stock.StockInfo),
	}

	client := redis.RedisClient
	pipe := client.Pipeline()

	// Setup commands
	stockCmds := make(map[uint32]*redis.StringCmd)
	frozenCmds := make(map[uint32]*redis.StringCmd)

	for _, item := range items {
		productId := item.ProductId
		stockKey := fmt.Sprintf("prod_{%d}_stock", productId)
		frozenKey := fmt.Sprintf("prod_{%d}_frozen", productId)

		stockCmds[productId] = pipe.Get(ctx, stockKey)
		frozenCmds[productId] = pipe.Get(ctx, frozenKey)
	}

	// Execute pipeline
	_, err := pipe.Exec(ctx)
	if err != nil {
		klog.Warnf("Failed to get stock from Redis: %v", err)
		return resp
	}

	// Process results
	for _, item := range items {
		productId := item.ProductId
		var stockVal, frozenVal int64 = 0, 0

		stockStr, err := stockCmds[productId].Result()
		if err == nil {
			stockVal, _ = strconv.ParseInt(stockStr, 10, 32)
		}

		frozenStr, err := frozenCmds[productId].Result()
		if err == nil {
			frozenVal, _ = strconv.ParseInt(frozenStr, 10, 32)
		}

		resp.Stocks[productId] = &stock.StockInfo{
			ProductId:   productId,
			Stock:       int32(stockVal),
			FrozenStock: int32(frozenVal),
		}
	}

	return resp
}
