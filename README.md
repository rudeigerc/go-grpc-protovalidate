# go-grpc-protovalidate

> [!CAUTION]
> WORK IN PROGRESS

A [gRPC Go](https://github.com/grpc/grpc-go) middleware for validating messages via [bufbuild/protovalidate](https://github.com/bufbuild/protovalidate) following Google's API Improvement Proposals [AIP-193](https://google.aip.dev/193).

## Installation

```shell
go get github.com/rudeigerc/go-grpc-protovalidate
```

## Usage

```go
package main

import (
	"log"

	"github.com/bufbuild/protovalidate-go"
	go_grpc_protovalidate "github.com/rudeigerc/go-grpc-protovalidate"
	"google.golang.org/grpc"
)

func main() {
	validator, err := protovalidate.New()
	if err != nil {
		log.Fatalf("failed to create validator: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			go_grpc_protovalidate.UnaryServerInterceptor(go_grpc_protovalidate.WithValidator(validator)),
		),
		grpc.ChainStreamInterceptor(
			go_grpc_protovalidate.StreamServerInterceptor(go_grpc_protovalidate.WithValidator(validator)),
		),
	)
}
```

## References

- [bufbuild/protovalidate](https://github.com/bufbuild/protovalidate)
- [grpc-ecosystem/go-grpc-middleware/interceptors/protovalidate](https://github.com/grpc-ecosystem/go-grpc-middleware/tree/main/interceptors/protovalidate)
- [Errors  |  Cloud API Design Guide  |  Google Cloud](https://cloud.google.com/apis/design/errors)
- [AIP-193: Errors](https://google.aip.dev/193)

## License

Apache 2.0
