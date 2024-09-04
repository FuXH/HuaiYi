module chat_service

go 1.21

replace github.com/FuXH/HuaiYi/chat_service => ./proto

require (
	github.com/tencent/vectordatabase-sdk-go v1.3.2
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common v1.0.967
	github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/hunyuan v1.0.967
	google.golang.org/grpc v1.65.0
)

require (
	github.com/clbanning/mxj v1.8.4 // indirect
	github.com/google/go-querystring v1.0.0 // indirect
	github.com/mitchellh/mapstructure v1.4.3 // indirect
	github.com/mozillazg/go-httpheader v0.2.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/tencentyun/cos-go-sdk-v5 v0.7.54 // indirect
	golang.org/x/net v0.25.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
	golang.org/x/text v0.15.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240528184218-531527333157 // indirect
	google.golang.org/protobuf v1.34.2 // indirect
)
