#!/bin/bash
# 设置变量
IMAGE_NAME="cloud-storage"
VERSION="latest"
CONTAINER_NAME="cloud-storage-container"

# 切换到项目根目录
cd "$(dirname "$0")/.."

# 构建 Docker 镜像
echo "Building Docker image..."
docker build -t ${IMAGE_NAME}:${VERSION} -f cmd/Dockerfile .

# 检查构建是否成功
if [ $? -eq 0 ]; then
    echo "Docker image built successfully!"
    echo "Image: ${IMAGE_NAME}:${VERSION}"
else
    echo "Error: Docker build failed!"
    exit 1
fi

# 运行新容器
echo "Starting new container..."
docker run -d \
    --name ${CONTAINER_NAME} \
    -p 8080:8080 \
    -p 6066:6066 \
    -v $(pwd)/logs:/app/cloud-storage/logs \
    -v $(pwd)/internal/initialize/config/config.yml:/app/cloud-storage/internal/initialize/config/config.yml \
    -v $(pwd)/internal/dist:/app/cloud-storage/internal/dist \
    ${IMAGE_NAME}:${VERSION}

# 检查容器是否成功启动
if [ $? -eq 0 ]; then
    echo "Container started successfully!"
    echo "Container name: ${CONTAINER_NAME}"
else
    echo "Error: Failed to start container!"
    exit 1
fi
