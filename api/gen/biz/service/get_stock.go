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

type GetStockService struct {
	ctx context.Context
}

// NewGetStockService new GetStockService
func NewGetStockService(ctx context.Context) *GetStockService {
	return &GetStockService{ctx: ctx}
}

// Run create note info
func (s *GetStockService) Run(req *stock.GetStockReq) (resp *stock.StockResp, err error) {
	// Initialize response
	resp = &stock.StockResp{
		Success: true,
		Stocks:  make(map[uint32]*stock.StockInfo),
	}

	// If no product IDs provided, return empty response
	if len(req.ProductIds) == 0 {
		return resp, nil
	}

	// Check cache first
	redisClient := redis.RedisClient
	pipe := redisClient.Pipeline()

	// Setup commands for getting stock and frozen stock for each product
	stockCmds := make(map[uint32]*redis.StringCmd)
	frozenCmds := make(map[uint32]*redis.StringCmd)

	for _, productID := range req.ProductIds {
		stockKey := fmt.Sprintf("prod_{%d}_stock", productID)
		frozenKey := fmt.Sprintf("prod_{%d}_frozen", productID)

		stockCmds[productID] = pipe.Get(s.ctx, stockKey)
		frozenCmds[productID] = pipe.Get(s.ctx, frozenKey)
	}

	// Execute pipeline
	_, err = pipe.Exec(s.ctx)
	if err != nil {
		klog.Warnf("Failed to get stock from Redis: %v", err)
	}

	// Process results and identify missing products
	missingProducts := make([]uint32, 0)
	for _, productID := range req.ProductIds {
		stockValue, stockErr := stockCmds[productID].Result()
		frozenValue, frozenErr := frozenCmds[productID].Result()

		if stockErr != nil || frozenErr != nil {
			missingProducts = append(missingProducts, productID)
			continue
		}

		stock, _ := strconv.ParseInt(stockValue, 10, 32)
		frozen, _ := strconv.ParseInt(frozenValue, 10, 32)

		resp.Stocks[productID] = &stock.StockInfo{
			ProductId:   productID,
			Stock:       int32(stock),
			FrozenStock: int32(frozen),
		}
	}

	// If we have missing products, fetch from database
	if len(missingProducts) > 0 {
		stocks, dbErr := model.BatchGetStockByProductIDs(mysql.DB, s.ctx, missingProducts)
		if dbErr != nil {
			klog.Errorf("Failed to get stock from database: %v", dbErr)
			// Continue with what we have from cache
		} else {
			// Update cache and response
			newPipe := redisClient.Pipeline()

			for _, stockItem := range stocks {
				productID := stockItem.ProductID

				// Add to response
				resp.Stocks[productID] = &stock.StockInfo{
					ProductId:   productID,
					Stock:       stockItem.Stock,
					FrozenStock: stockItem.FrozenStock,
				}

				// Update cache
				newPipe.Set(s.ctx, fmt.Sprintf("prod_{%d}_stock", productID), stockItem.Stock, 0)
				newPipe.Set(s.ctx, fmt.Sprintf("prod_{%d}_frozen", productID), stockItem.FrozenStock, 0)
			}

			_, cacheErr := newPipe.Exec(s.ctx)
			if cacheErr != nil {
				klog.Warnf("Failed to update stock cache: %v", cacheErr)
			}
		}
	}

	return resp, nil
}
