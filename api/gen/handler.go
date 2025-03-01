package main

import (
	"DouyinMerchant/api/gen/biz/service"
	"DouyinMerchant/api/gen/kitex_gen/cart"
	auth "DouyinMerchant/api/gen/kitex_gen/douyin_merchant/aut
	"DouyinMerchant/api/gen/kitex_gen/douyin_merchant/product"
"
	"DouyinMerchant/api/gen/kitex_gen/douyin_merchant/sto
k"
	"DouyinMerchant/api/gen/kitex_gen/douyin_merchant/us
	"DouyinMerchant/api/gen/kitex_gen/douyin_merchant/user"
	"context"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}
type UserServiceImpl struct{}
type ProductServiceImpl struct{}

// DeliverTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeliverTokenByRPC(ctx context.Context, req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	resp, err = service.NewDeliverTokenByRPCService(ctx).Run(req)

	return resp, err
}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	resp, err = service.NewVerifyTokenByRPCService(ctx).Run(req)

	return resp, err
}

// Register implements the UserServiceImpl interface.
func (s *UserServiceImpl) Register(ctx context.Context, req *user.RegisterReq) (resp *user.RegisterResp, err error) {
	resp, err = service.NewRegisterService(ctx).Run(req)

	return resp, err
}

// Login implements the UserServiceImpl interface.
func (s *UserServiceImpl) Login(ctx context.Context, req *user.LoginReq) (resp *user.LoginResp, err error) {
	resp, err = service.NewLoginService(ctx).Run(req)

	return resp, err
}

// ListProducts implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) ListProducts(ctx context.Context, req *product.ListProductsReq) (resp *product.ListProductsResp, err error) {
	resp, err = service.NewListProductsService(ctx).Run(req)

	return resp, err
}

// GetProduct implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) GetProduct(ctx context.Context, req *product.GetProductReq) (resp *product.GetProductResp, err error) {
	resp, err = service.NewGetProductService(ctx).Run(req)

	return resp, err
}

// SearchProducts implements the ProductServiceImpl interface.
func (s *ProductServiceImpl) SearchProducts(ctx context.Context, req *product.SearchProductsReq) (resp *product.SearchProductsResp, err error) {
	resp, err = service.NewSearchProductsService(ctx).Run(req)

	return resp, err
}

// AddItem implements the CartServiceImpl interface.
func (s *CartServiceImpl) AddItem(ctx context.Context, req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	resp, err = service.NewAddItemService(ctx).Run(req)

	return resp, err
}

// GetCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) GetCart(ctx context.Context, req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	resp, err = service.NewGetCartService(ctx).Run(req)

	return resp, err
}

// EmptyCart implements the CartServiceImpl interface.
func (s *CartServiceImpl) EmptyCart(ctx context.Context, req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	resp, err = service.NewEmptyCartService(ctx).Run(req)

	return resp, err
}

// PlaceOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) PlaceOrder(ctx context.Context, req *order.PlaceOrderReq) (resp *order.PlaceOrderResp, err error) {
	resp, err = service.NewPlaceOrderService(ctx).Run(req)

	return resp, err
}

// ListOrder implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) ListOrder(ctx context.Context, req *order.ListOrderReq) (resp *order.ListOrderResp, err error) {
	resp, err = service.NewListOrderService(ctx).Run(req)

	return resp, err
}

// MarkOrderPaid implements the OrderServiceImpl interface.
func (s *OrderServiceImpl) MarkOrderPaid(ctx context.Context, req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {
	resp, err = service.NewMarkOrderPaidService(ctx).Run(req)

	return resp, err
}

// DecreaseStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) DecreaseStock(ctx context.Context, req *stock.DecreaseStockReq) (resp *stock.StockResp, err error) {
	resp, err = service.NewDecreaseStockService(ctx).Run(req)

	return resp, err
}

// GetStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) GetStock(ctx context.Context, req *stock.GetStockReq) (resp *stock.StockResp, err error) {
	resp, err = service.NewGetStockService(ctx).Run(req)

	return resp, err
}

// PreReduceStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) PreReduceStock(ctx context.Context, req *stock.PreReduceStockReq) (resp *stock.StockResp, err error) {
	resp, err = service.NewPreReduceStockService(ctx).Run(req)

	return resp, err
}

// RollbackStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) RollbackStock(ctx context.Context, req *stock.RollbackStockReq) (resp *stock.StockResp, err error) {
	resp, err = service.NewRollbackStockService(ctx).Run(req)

	return resp, err
}

// BatchAddStock implements the StockServiceImpl interface.
func (s *StockServiceImpl) BatchAddStock(ctx context.Context, req *stock.BatchAddStockReq) (resp *stock.StockResp, err error) {
	resp, err = service.NewBatchAddStockService(ctx).Run(req)

	return resp, err
}
