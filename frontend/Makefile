# 项目相关配置
PROJECT_NAME:=ld-forge
BUILD_NUMBER:=loacl
PORT:=8840

# 私有仓库配置
REGISTRY:=1.117.192.82:8666
DOCKER_USERNAME:=admin
DOCKER_PASSWORD:=yskj2407

# 路径配置
NGINX_CONF_DIR:=nginx
DOCKERFILE:=Dockerfile

IMAGE_NAME := $(REGISTRY)/$(PROJECT_NAME)-fe:$(BUILD_NUMBER)
CONTAINER_NAME := $(PROJECT_NAME)-fe

.PHONY: all clean docker-login build run

all: clean build run

clean:
	@echo "停止并删除容器$(CONTAINER_NAME)..."
	docker stop $(CONTAINER_NAME) || true
	docker rm $(CONTAINER_NAME) || true

# 登录私有仓库
docker-login:
	@echo "登录仓库$(REGISTRY)"
	echo $(DOCKER_PASSWORD) | docker login $(REGISTRY) -u $(DOCKER_USERNAME) --password-stdin
	@echo "登录成功"

build: docker-login
	@echo "开始构建镜像$(IMAGE_NAME)"
	docker build -t $(IMAGE_NAME) -f $(DOCKERFILE) .
	@echo "构建完成"

run:
	@echo "运行容器$(CONTAINER_NAME)..."
	docker run -it --rm \
		-p $(PORT):$(PORT) \
		-v $(PWD)/$(NGINX_CONF_DIR)/nginx.conf:/etc/nginx/nginx.conf \
		-v $(PWD)/$(NGINX_CONF_DIR)/fe_dev.conf:/etc/nginx/conf.d/default.conf \
		--name $(CONTAINER_NAME) \
		$(IMAGE_NAME)
	@echo "容器已经启动:http://localhost:$(PORT)/"