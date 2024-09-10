.PHONY: .proto

generate:
	protoc \
		--proto_path=api/ \
		--go_out=internal/pkg/pb \
		--go-grpc_out=internal/pkg/pb \
		api/*.proto

