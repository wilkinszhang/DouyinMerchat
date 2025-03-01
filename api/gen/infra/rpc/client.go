package rpc

import (
	"DouyinMerchant/api/gen/conf"
	"DouyinMerchant/api/gen/kitex_gen/douyin_merchant/product/productservice"
	"DouyinMerchant/api/gen/utils"
	"DouyinMerchant/common/clientsuite"
	"github.com/cloudwego/kitex/client"
	"sync"
)

var (
	ProductClient productservice.Client
	once          sync.Once
	err           error
	registryAddr  string
	serviceName   string
)

func InitClient() {
	once.Do(func() {
		registryAddr = conf.GetConf().Registry.RegistryAddress[0]
		serviceName = conf.GetConf().Kitex.Service
		initProductClient()
	})
}

func initProductClient() {
	opts := []client.Option{
		client.WithSuite(clientsuite.CommonGrpcClientSuite{
			RegistryAddr:       registryAddr,
			CurrentServiceName: serviceName,
		}),
	}

	ProductClient, err = productservice.NewClient("product", opts...)
	utils.MustHandleError(err)
}
