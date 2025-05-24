#!/bin/bash

# 设置变量
IMAGE_NAME="cloud-storage-backend"
VERSION="latest"
CONFIG_FILE="../internal/initialize/config/config.yml"
CONFIG_DIR="../internal/initialize/config"

# 确保配置文件存在
if [ ! -f "$CONFIG_FILE" ]; then
    echo "Error: Config file not found at $CONFIG_FILE"
    exit 1
fi

# 确保配置目录存在
if [ ! -d "$CONFIG_DIR" ]; then
    echo "Error: Config directory not found at $CONFIG_DIR"
    exit 1
fi

# 构建 Docker 镜像
echo "Building Docker image..."
docker build -t ${IMAGE_NAME}:${VERSION} .

# 检查构建是否成功
if [ $? -eq 0 ]; then
    echo "Docker image built successfully!"
    echo "Image: ${IMAGE_NAME}:${VERSION}"
else
    echo "Error: Docker build failed!"
    exit 1
fi 