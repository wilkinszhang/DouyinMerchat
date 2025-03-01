package service

import (
	"DouyinMerchant/api/gen/biz/dal"
	product "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/product"
	"context"
	"github.com/joho/godotenv"
	"testing"
)

func TestListProducts_Run(t *testing.T) {
	err := godotenv.Load()
	dal.Init()
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	// init req and assert value

	// todo: edit your unit test
	//进行单测，以及jmeter压测
	t.Run("Empty Request - Default Values", func(t *testing.T) {
		req := &product.ListProductsReq{}
		s := NewListProductsService(ctx)
		s.Run(req)
	})

	t.Run("Invalid Page Number", func(t *testing.T) {
		req := &product.ListProductsReq{
			Page:     0,
			PageSize: 10,
		}
		s := NewListProductsService(ctx)
		_, err := s.Run(req)

		if err == nil {
			t.Error("Expected error for invalid page number, got nil")
		}
	})

	t.Run("Invalid Page Size", func(t *testing.T) {
		req := &product.ListProductsReq{
			Page:     1,
			PageSize: 0,
		}
		s := NewListProductsService(ctx)
		_, err := s.Run(req)

		if err == nil {
			t.Error("Expected error for invalid page size, got nil")
		}
	})

	t.Run("Filter By Category", func(t *testing.T) {
		req := &product.ListProductsReq{
			Page:         1,
			PageSize:     10,
			CategoryName: "数码电器",
		}
		s := NewListProductsService(ctx)
		resp, err := s.Run(req)

		if err != nil {
			t.Errorf("Unexpected error: %v", err)
		}
		// Add assertions for filtered results
		for _, product := range resp.Products {
			foundCategory := false
			for _, category := range product.Categories {
				if category == req.CategoryName {
					foundCategory = true
					break
				}
			}
			if !foundCategory {
				t.Errorf("Expected product %s to have category %s, but category not found",
					product.Name, req.CategoryName)
			}
		}
	})
}
