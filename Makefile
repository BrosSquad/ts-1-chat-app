PROTO_FILES ?= $(shell find proto -iname "*.proto")
GOPATH ?= $(HOME)/go
LOG_LEVEL ?= trace
PROTO_GENERATE_PATH = backend/services/pb

.PHONY: build
build: clean
	@cd backend/ && go build -o ../bin/server ./cmd/server/main.go
	@cp backend/config.example.yml ./bin/config.yml
	@mkdir ./bin/logs
	@mkdir ./bin/db
	@echo './server -db ./db/database.sqlite -logs ./logs -config . -file -console -level info' >> ./bin/start
	@chmod +x ./bin/start

config.yml:
	@cp backend/config.example.yml ./config.yml

run: config.yml
	@cd backend && go run cmd/server/main.go \
		-logs ../logs \
		-config .. \
		-level $(LOG_LEVEL) \
		-file \
		-console

.PHONY: protoc-go
protoc-go:
	@protoc -I proto -I proto-3rd-party \
	 	--go_out=$(PROTO_GENERATE_PATH) \
		--go-grpc_out=$(PROTO_GENERATE_PATH) \
		--go-tag_out=paths=source_relative:$(PROTO_GENERATE_PATH) \
		--go-grpc_opt=paths=source_relative \
		$(shell find proto -iname "*.proto")
	@rm -rf $(PROTO_GENERATE_PATH)/github.com

.PHONY: protoc-dart
protoc-dart:
	@protoc -I proto -I proto-3rd-party --dart_out=grpc:chat_app/lib/proto \
		$(PROTO_FILES)

.PHONY: clean
clean:
	@rm -rf bin/