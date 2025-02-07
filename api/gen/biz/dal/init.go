package dal

import (
	"DouyinMerchant/api/gen/biz/dal/mysql"
	"DouyinMerchant/api/gen/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
