.PHONY: build
build:
	@go build -o bin/server backend/cmd/server/main.go


.PHONY: protoc-go
protoc-go:
	protoc -I proto --go_out=backend/services/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=backend/services/pb \
		--go-grpc_opt=paths=source_relative $(shell find proto -iname "*.proto")

.PHONY: protoc-dart
protoc-dart:
	protoc -I proto --dart_out=grpc:chat_app/lib/proto \
		$(shell find proto -iname "*.proto")
