package service

import (
	"DouyinMerchant/api/gen/biz/dal/mysql"
	"DouyinMerchant/api/gen/biz/model"
	product "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/product"
	"context"
	"errors"
	"strings"
)

type ListProductsService struct {
	ctx context.Context
} // NewListProductsService new ListProductsService
func NewListProductsService(ctx context.Context) *ListProductsService {
	return &ListProductsService{ctx: ctx}
}

// Run create note info
func (s *ListProductsService) Run(req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	// Finish your business logic.
	//先校验参数合法
	if err := validateRequest(req); err != nil {
		return nil, err
	}
	//传入类别名字，返回所有category列
	products, err := model.GetProductsByCategoryName(mysql.DB, s.ctx, req.CategoryName, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}
	resp = &product.ListProductsResp{Products: make([]*product.Product, 0)}
	//遍历所有商品
	for _, p := range products {
		resp.Products = append(resp.Products, &product.Product{
			Id:          uint32(p.ID),
			Name:        p.Name,
			Description: p.Description,
			Picture:     p.Picture,
			Price:       p.Price,
			Categories:  strings.Split(p.Categories, ","), // 将逗号分隔的分类字符串转为切片
		})
	}
	return resp, nil
}

func validateRequest(req *product.ListProductsReq) error {
	if req.Page < 1 {
		return errors.New("page must be greater than 0")
	}
	if req.PageSize < 1 {
		return errors.New("page_size must be greater than 0")
	}
	if strings.TrimSpace(req.CategoryName) == "" {
		return errors.New("category_name cannot be empty")
	}
	return nil
}
