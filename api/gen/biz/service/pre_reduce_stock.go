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

type PreReduceStockService struct {
	ctx context.Context
}

// NewPreReduceStockService new PreReduceStockService
func NewPreReduceStockService(ctx context.Context) *PreReduceStockService {
	return &PreReduceStockService{ctx: ctx}
}

// Run create note info
func (s *PreReduceStockService) Run(req *stock.PreReduceStockReq) (resp *stock.StockResp, err error) {
	// Initialize response
	resp = &stock.StockResp{
		Success: true,
		Stocks:  make(map[uint32]*stock.StockInfo),
	}

	// Redis client
	client := redis.RedisClient
	ctx := s.ctx

	// Redis transaction to pre-reduce stock
	pipe := client.TxPipeline()

	stockMap := make(map[uint32]int32)
	for _, item := range req.Items {
		// Use hash tag to ensure keys are in the same Redis slot
		stockKey := fmt.Sprintf("prod_{%d}_stock", item.ProductId)
		frozenKey := fmt.Sprintf("prod_{%d}_frozen", item.ProductId)

		// Get current stock value
		currentStock, err := client.Get(ctx, stockKey).Result()
		if err != nil {
			// If key doesn't exist, try to fetch from database
			stockData, dbErr := model.GetStockByProductID(mysql.DB, ctx, item.ProductId)
			if dbErr != nil {
				resp.Success = false
				resp.Message = fmt.Sprintf("Failed to get stock for product %d", item.ProductId)
				return resp, nil
			}

			// Initialize Redis with database values
			pipe.Set(ctx, stockKey, stockData.Stock, 0)
			pipe.Set(ctx, frozenKey, stockData.FrozenStock, 0)

			currentStock = fmt.Sprintf("%d", stockData.Stock)
		}

		// Parse stock value
		stockNum, err := strconv.Atoi(currentStock)
		if err != nil {
			resp.Success = false
			resp.Message = fmt.Sprintf("Invalid stock value for product %d", item.ProductId)
			return resp, nil
		}

		// Check if enough stock is available
		if stockNum < int(item.Num) {
			resp.Success = false
			resp.Message = fmt.Sprintf("Insufficient stock for product %d", item.ProductId)
			return resp, nil
		}

		// Track stock changes for database update
		stockMap[item.ProductId] = item.Num

		// Add commands to pipeline
		pipe.DecrBy(ctx, stockKey, int64(item.Num))
		pipe.IncrBy(ctx, frozenKey, int64(item.Num))
	}

	// Execute Redis transaction
	_, err = pipe.Exec(ctx)
	if err != nil {
		resp.Success = false
		resp.Message = "Concurrent modification, please try again"
		return resp, nil
	}

	// Update database asynchronously
	go s.asyncUpdateDB(stockMap)

	// Get updated stock for response
	for productId, _ := range stockMap {
		stockValue, err := client.Get(ctx, fmt.Sprintf("prod_{%d}_stock", productId)).Result()
		frozenValue, err2 := client.Get(ctx, fmt.Sprintf("prod_{%d}_frozen", productId)).Result()

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
}

func (s *PreReduceStockService) asyncUpdateDB(stockMap map[uint32]int32) {
	// Create new context for the goroutine
	ctx := context.Background()

	tx := mysql.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	for productId, num := range stockMap {
		err := model.PreReduceStock(tx, ctx, productId, num)
		if err != nil {
			klog.Errorf("Failed to update stock in database: %v", err)
			tx.Rollback()

			// Attempt to rollback Redis changes
			s.rollbackRedis(productId, num)
			return
		}
	}

	if err := tx.Commit().Error; err != nil {
		klog.Errorf("Failed to commit stock transaction: %v", err)

		// Rollback Redis for all products
		for productId, num := range stockMap {
			s.rollbackRedis(productId, num)
		}
	}
}

func (s *PreReduceStockService) rollbackRedis(productId uint32, num int32) {
	ctx := context.Background()
	client := redis.RedisClient

	stockKey := fmt.Sprintf("prod_{%d}_stock", productId)
	frozenKey := fmt.Sprintf("prod_{%d}_frozen", productId)

	pipe := client.Pipeline()
	pipe.IncrBy(ctx, stockKey, int64(num))
	pipe.DecrBy(ctx, frozenKey, int64(num))

	_, err := pipe.Exec(ctx)
	if err != nil {
		klog.Errorf("Failed to rollback Redis for product %d: %v", productId, err)
	}
}

func getRedisStock(id uint32) int32 {
	ctx := context.Background()
	stockKey := fmt.Sprintf("prod_{%d}_stock", id)

	val, err := redis.RedisClient.Get(ctx, stockKey).Result()
	if err != nil {
		return 0
	}

	stock, _ := strconv.ParseInt(val, 10, 32)
	return int32(stock)
}

func getRedisFrozen(id uint32) int32 {
	ctx := context.Background()
	frozenKey := fmt.Sprintf("prod_{%d}_frozen", id)

	val, err := redis.RedisClient.Get(ctx, frozenKey).Result()
	if err != nil {
		return 0
	}

	frozen, _ := strconv.ParseInt(val, 10, 32)
	return int32(frozen)
}
