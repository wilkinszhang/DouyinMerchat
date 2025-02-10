.PHONY: gen-auth
gen-auth:
	@cd api/gen && cwgo server -I ../proto --module DouyinMerchant --service AuthService --idl ../proto/auth.proto
gen-user:
	@cd api/gen && cwgo server -I ../proto --module DouyinMerchant --service user_service --idl ../proto/user.proto
