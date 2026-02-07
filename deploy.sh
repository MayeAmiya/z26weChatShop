#!/bin/bash

# Z26 WeChat Shop 一键部署脚本 (WSL)
# 用于在 WSL 环境中快速部署前后端服务

set -e  # 遇到错误立即退出

echo "🚀 开始部署 Z26 WeChat Shop..."

# 检查 Docker 是否运行
if ! docker info > /dev/null 2>&1; then
    echo "❌ Docker 未运行，请先启动 Docker Desktop"
    exit 1
fi

# 检查 docker-compose 是否可用
if ! command -v docker-compose > /dev/null 2>&1; then
    echo "❌ docker-compose 未安装"
    exit 1
fi

echo "📦 构建后端镜像..."
docker build -t develop-backend:latest ./z26b-backend

echo "🎨 构建前端镜像..."
docker build -t develop-frontend:latest ./z26a

echo "🔄 启动服务..."
docker-compose -f docker-compose.deploy.yml up -d

echo "⏳ 等待服务启动..."
sleep 10

echo "📊 检查服务状态..."
docker-compose -f docker-compose.deploy.yml ps

echo ""
echo "✅ 部署完成！"
echo "🌐 前端地址: http://localhost"
echo "🔗 后端 API: http://localhost:8080"
echo "📦 MinIO 控制台: http://localhost:9001 (admin/admin123456)"
echo ""
echo "停止服务: docker-compose -f docker-compose.deploy.yml down"
echo "查看日志: docker-compose -f docker-compose.deploy.yml logs -f"