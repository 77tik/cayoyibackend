.PHONY: env genapi gensql clean lint run local-dev local-deps local-up local-down local-restart local-logs

# 工具链版本
gorm_gentool := gorm.io/gen/tools/gentool@v0.0.1
golangci-lint := github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.1.2
goctl := github.com/zeromicro/go-zero/tools/goctl@v1.8.4

# 本地开发配置
LOCAL_PORT ?= 8849
LOCAL_CONFIG ?= ./etc/ldhydropower-api-local.yaml
DOCKER_NETWORK ?= ld-hydropower-network
DB_CONTAINER ?= ld-hydropower-db
DB_PORT ?= 3306
DB_USER ?= root
DB_PASSWORD ?= root
DB_NAME ?= ldhydropowerdb
DB_VERSION ?= mysql:8.0

env:
#	@echo "=========install xgo, the All-in-one go testing library========="
#	GOPROXY=https://goproxy.cn/,direct go install github.com/xhd2015/xgo/cmd/xgo@latest

gensql:
	go run $(gorm_gentool) -dsn "$(DB_USER):$(DB_PASSWORD)@tcp(localhost:$(DB_PORT))/$(DB_NAME)?charset=utf8mb4&parseTime=True&loc=Asia%2FShanghai" \
		-outPath "./internal/dao/query" -fieldNullable -fieldSignable

genapi:
	go run $(goctl) api format --dir ./api
	go run $(goctl) api go --home ./.goctltmpl --api ./api/main.api --dir ./ --style go_zero --type-group
	goctl api plugin -plugin goctl-swagger="swagger -filename swagger.json" -api ./api/main.api -dir ./internal/handler/swagger/

lint:
	go run $(golangci-lint) run ./...

# 本地开发使用的运行命令
run: genapi
	go run . -f $(LOCAL_CONFIG)

clean:
	go clean -i .

# 本地开发环境相关目标
local-deps:
	# 创建网络
	docker network create $(DOCKER_NETWORK) || true
	# 启动数据库容器
	docker run -d \
		--name $(DB_CONTAINER) \
		--network $(DOCKER_NETWORK) \
		-p $(DB_PORT):3306 \
		-e MYSQL_ROOT_PASSWORD=$(DB_PASSWORD) \
		-e MYSQL_DATABASE=$(DB_NAME) \
		$(DB_VERSION)
	# 等待数据库启动
	@echo "Waiting for database to start..."
	@sleep 10

local-up: local-deps
	# 构建并启动应用容器
	docker build -t ld-hydropower-be:local .
	docker run -d \
		--name ld-hydropower-be \
		--network $(DOCKER_NETWORK) \
		-p $(LOCAL_PORT):$(LOCAL_PORT) \
		-v $(PWD)/etc:/app/etc \
		ld-hydropower-be:local -f /app/etc/ldhydropower-api-local.yaml

local-down:
	docker stop ld-hydropower-be $(DB_CONTAINER) || true
	docker rm ld-hydropower-be $(DB_CONTAINER) || true
	docker network rm $(DOCKER_NETWORK) || true

local-restart: local-down local-up

local-logs:
	docker logs -f ld-hydropower-be

