字节青训营 抖音电商

- [x] 为每个服务设置 protobuf 文件
- [ ] 实现身份验证中间件
- [ ] 创建数据库模式
- [ ] 逐个实现服务
- [ ] 设置服务发现 (etcd/consul)
- [ ] 添加可观察性 (日志记录/监控)

问题记录：api/gen/cmd/client/client.go测试报错，认证服务已经写好了，先不管他，调用RPC服务报错：
GOROOT=E:\go1.23.3 #gosetup
GOPATH=E:\go #gosetup
E:\go1.23.3\bin\go.exe build -o E:\Users\zhangweijian\AppData\Local\JetBrains\Goland\tmp\GoLand\___go_build_handler_go_main_go.exe E:\DouyinMerchant\api\gen\handler.go E:\DouyinMerchant\api\gen\main.go #gosetup
E:\Users\zhangweijian\AppData\Local\JetBrains\Goland\tmp\GoLand\___go_build_handler_go_main_go.exe #gosetup
&{Env:test Kitex:{Service:AuthService Address::8888 LogLevel:info LogFileName:log/kitex.log LogMaxSize:10 LogMaxBackups:50 LogMaxAge:3} MySQL:{DSN:%s:%s@tcp(%s:3310)/%s?charset=utf8mb4&parseTime=True&loc=Local} Redis:{Address:127.0.0.1:6379 Username: Password: DB:0} Registry:{RegistryAddress:[127.0.0.1:2379] Username: Password:}}

2025/02/07 16:24:14 E:/DouyinMerchant/api/gen/biz/dal/mysql/init.go:29
[0.524ms] [rows:0] select version()
&gorm.DB{Config:(*gorm.Config)(0xc001d21ef0), Error:error(nil), RowsAffected:0, Statement:(*gorm.Statement)(0xc001d616c0), clone:0}panic: close of closed channel

goroutine 41 [running]:
github.com/cloudwego/kitex/pkg/remote/trans/nphttp2/grpc.(*http2Server).HandleStreams(0xc001cc6000, 0xc001ca4280, 0x24f18b0)
E:/go/pkg/mod/github.com/cloudwego/kitex@v0.12.1/pkg/remote/trans/nphttp2/grpc/http2_server.go:447 +0xc4d
github.com/cloudwego/kitex/pkg/remote/trans/nphttp2.(*svrTransHandler).OnRead(0xc001d1b1c0, {0x2778cf0, 0xc001c923c0}, {0x277de48, 0xc001c8c080})
E:/go/pkg/mod/github.com/cloudwego/kitex@v0.12.1/pkg/remote/trans/nphttp2/server_handler.go:133 +0x104
github.com/cloudwego/kitex/pkg/remote/trans/detection.(*svrTransHandler).OnRead(0xc001d62d50, {0x2778cf0, 0xc001c92150}, {0x277de48, 0xc001c8c080})
E:/go/pkg/mod/github.com/cloudwego/kitex@v0.12.1/pkg/remote/trans/detection/server_handler.go:95 +0xb7
github.com/cloudwego/kitex/pkg/remote.(*TransPipeline).OnRead(0xc001d1b200, {0x2778cf0?, 0xc001c92150?}, {0x277de48, 0xc001c8c080})
E:/go/pkg/mod/github.com/cloudwego/kitex@v0.12.1/pkg/remote/trans_pipeline.go:129 +0xbb
github.com/cloudwego/kitex/pkg/remote/trans/gonet.(*transServer).BootstrapServer.func1()
E:/go/pkg/mod/github.com/cloudwego/kitex@v0.12.1/pkg/remote/trans/gonet/trans_server.go:101 +0x223
created by github.com/cloudwego/kitex/pkg/remote/trans/gonet.(*transServer).BootstrapServer in goroutine 39
E:/go/pkg/mod/github.com/cloudwego/kitex@v0.12.1/pkg/remote/trans/gonet/trans_server.go:85 +0x88

Process finished with the exit code 2

用户服务：Register ok，Login ok

