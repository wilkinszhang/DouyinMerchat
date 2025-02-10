package main

import (
	"DouyinMerchant/api/gen/conf"
	"DouyinMerchant/api/gen/kitex_gen/douyin_merchant/user"
	"DouyinMerchant/api/gen/kitex_gen/douyin_merchant/user/userservice"
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	consul "github.com/kitex-contrib/registry-consul"
	"log"
)

var (
	cli userservice.Client
)

func main() {
	r, err := consul.NewConsulResolver(conf.GetConf().Registry.RegistryAddress[0])
	if err != nil {
		panic(err)
	}
	c, err := userservice.NewClient("user_service", client.WithResolver(r),
		client.WithTransportProtocol(transport.GRPC),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		//client.WithFailureRetry(retry.NewFailurePolicy()),       //增加重试
		//client.WithCircuitBreaker(circuitbreak.NewCBSuite(nil)), //增加熔断
		//client.WithRPCTimeout(3*time.Second),                    //超时
	)
	cli = c
	if err != nil {
		panic(err)
	}
	hz := server.New(
		server.WithHostPorts("localhost:8181"),
		//server.WithMaxRequestBodySize(4*1024*1024), //限制body大小
		//server.WithIdleTimeout(300*time.Second),    //限制空闲超时
	)

	hz.POST("/user_service", Handler)

	if err := hz.Run(); err != nil {
		log.Fatal(err)
	}
	//res, err := c.Register(context.Background(), &user.RegisterReq{Email: "2658536235@qq.com", Password: "test", ConfirmPassword: "test"})
	//fmt.Printf("%v", res)
}

func Handler(ctx context.Context, c *app.RequestContext) {
	req := user.RegisterReq{}
	//req.Email = "22@bilibili.com"
	//req.Password = "123"
	//req.ConfirmPassword = "123"
	if err := c.Bind(&req); err != nil {
		c.String(400, err.Error())
		return
	}
	resp, err := cli.Register(context.Background(), &req)
	if err != nil {
		c.String(500, err.Error())
		log.Fatal(err)
	}

	c.String(200, resp.String())
}
