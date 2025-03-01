package service

import (
	"DouyinMerchant/api/gen/biz/dal/mysql"
	"DouyinMerchant/api/gen/biz/model"
	product "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/product"
	"context"
	"strings"
)

type SearchProductsService struct {
	ctx context.Context
} // NewSearchProductsService new SearchProductsService
func NewSearchProductsService(ctx context.Context) *SearchProductsService {
	return &SearchProductsService{ctx: ctx}
}

// Run create note info
func (s *SearchProductsService) Run(req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	// Finish your business logic.
	//这个接口应该用ElasticSearch实现最好。
	query := strings.TrimSpace(req.Query)
	if query == "" {
		return &product.SearchProductsResp{}, nil
	}

	// 调用 model 包中的查询方法
	productQuery := model.NewProductQuery(context.Background(), mysql.DB)
	results, err := productQuery.GetProductsByQuery(query)
	if err != nil {
		return nil, err
	}

	// 转换结果为响应格式
	resp = &product.SearchProductsResp{
		Results: make([]*product.Product, 0, len(results)),
	}

	for _, r := range results {
		categories := strings.Split(r.Categories, ",")
		if r.Categories == "" {
			categories = []string{}
		}

		resp.Results = append(resp.Results, &product.Product{
			Id:          uint32(int64(r.ID)),
			Name:        r.Name,
			Description: r.Description,
			Picture:     r.Picture,
			Price:       r.Price,
			Categories:  categories,
		})
	}

	return resp, nil
}
