package main

import (
	"DouyinMerchant/api/gen/biz/dal"
	"DouyinMerchant/api/gen/biz/service"
	"DouyinMerchant/api/gen/kitex_gen/douyin_merchant/product/productservice"
	"DouyinMerchant/mq"
	consul "github.com/kitex-contrib/registry-consul"
	"net"
	"time"

	"DouyinMerchant/api/gen/conf"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/joho/godotenv"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	opts := kitexInit()
	//增加限流
	//opts = append(opts, server.WithLimit(&limit.Option{
	//	MaxConnections: 1000,
	//	MaxQPS:         500,
	//}))

	dal.Init()

	// 初始化生产者
	if err := mq.InitProducer(); err != nil {
		panic(err)
	}
	defer mq.ShutdownProducer()
	go service.StartOrderConsumer()
	//[PlaceOrderService] --> (DB)
	//	↓
	//	[RocketMQ Producer] --延迟消息--> [RocketMQ Broker]
	//	↓
	//	[OrderConsumer] --> (检查订单状态) --> [更新订单状态]

	svr := productservice.NewServer(new(ProductServiceImpl), opts...)

	err = svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	// service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))

	r, err := consul.NewConsulRegister(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		klog.Fatal(err)
	}
	opts = append(opts, server.WithRegistry(r))

	// klog
	logger := kitexlogrus.NewLogger()
	klog.SetLogger(logger)
	klog.SetLevel(conf.LogLevel())
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.GetConf().Kitex.LogFileName,
			MaxSize:    conf.GetConf().Kitex.LogMaxSize,
			MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
			MaxAge:     conf.GetConf().Kitex.LogMaxAge,
		}),
		FlushInterval: time.Minute,
	}
	klog.SetOutput(asyncWriter)
	server.RegisterShutdownHook(func() {
		asyncWriter.Sync()
	})
	return
}
