.PHONY: build
build:
	@go build -o bin/server backend/main.go


.PHONY: protoc-go
protoc-go:
	cd proto/ && protoc -I. --go_out=../backend/services/pb \
		--go_opt=paths=source_relative \
		--go-grpc_out=../backend/services/pb \
		--go-grpc_opt=paths=source_relative \*.proto
