.PHONY: build
build:
	@cd backend/ && go build -o ../bin/server ./cmd/server/main.go
	cp backend/config.example.yml ./bin/config.yml
	mkdir ./bin/logs
	mkdir ./bin/db
	echo './server -db ./db/database.sqlite -logs ./logs -config . -file -console -level info' >> ./bin/start
	chmod +x ./bin/start

run:
	cd backend && go run cmd/server/main.go \
		-logs ../logs \
		-config .. \
		-level trace \
		-file \
		-console

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
