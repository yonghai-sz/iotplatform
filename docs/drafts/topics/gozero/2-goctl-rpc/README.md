
`goctl rpc template -o transform.proto`

`goctl rpc protoc transform.proto --go_out=./pb --go-grpc_out=./pb --zrpc_out=.`



`go install github.com/golang/mock/mockgen@v1.6.0`  
`mockgen -source=transformer/transformer.go -destination=transformer/transformer_mock.go -package=transformer`



要修改的地方：  
* etc/example.yaml
* internal/config/config.go
* internal/svc/servicecontext.go
* internal/logic/*.go


