#!/bin/bash

# 设置变量
IMAGE_NAME="cloud-storage-backend"
CONTAINER_NAME="cloud-storage-backend-container"

echo "🧹 开始清理 Docker 资源..."

# 停止并删除容器
echo "正在停止并删除容器 ${CONTAINER_NAME}..."
docker stop ${CONTAINER_NAME} 2>/dev/null || true
docker rm ${CONTAINER_NAME} 2>/dev/null || true

# 删除镜像
echo "正在删除镜像 ${IMAGE_NAME}..."
docker rmi ${IMAGE_NAME}:latest 2>/dev/null || true

# 清理未使用的镜像和容器
echo "正在清理未使用的 Docker 资源..."
docker system prune -f

echo "✅ 清理完成！"
echo "已删除："
echo "- 容器: ${CONTAINER_NAME}"
echo "- 镜像: ${IMAGE_NAME}:latest"
echo "- 未使用的 Docker 资源" 