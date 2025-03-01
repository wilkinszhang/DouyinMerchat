.PHONY: gen-auth
gen-auth:
	@cd api/gen && cwgo server -I ../proto --module DouyinMerchant --service AuthService --idl ../proto/auth.proto
gen-user:
	@cd api/gen && cwgo server -I ../proto --module DouyinMerchant --service user_service --idl ../proto/user.proto
gen-product:
	@cd api/gen && cwgo server -I ../proto --module DouyinMerchant --service product_service --idl ../proto/product.proto
gen-cart:
	@cd api/gen && cwgo server -I ../proto --module DouyinMerchant --service cart_service --idl ../proto/cart.proto
gen-order:
	@cd api/gen && cwgo server -I ../proto --module DouyinMerchant --service order_service --idl ../proto/order.proto
gen-stock:
	@cd api/gen && cwgo server -I ../proto --module DouyinMerchant --service stock_service --idl ../proto/stock.proto
gen-client:
	@cd api/gen && cwgo client --type RPC --service product_service --module DouyinMerchant  -I ../proto  --idl ../proto/product.proto