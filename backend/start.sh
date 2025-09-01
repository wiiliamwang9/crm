#!/bin/bash

echo "=== CRM客户管理系统启动脚本 ==="

# 检查端口是否被占用
if lsof -i:8080 > /dev/null 2>&1; then
    echo "警告: 端口8080已被占用，请先关闭占用该端口的进程"
    echo "可以使用命令: lsof -i:8080 查看占用端口的进程"
    exit 1
fi

# 检查web目录是否存在
if [ ! -d "../web" ]; then
    echo "错误: 找不到web目录，请确保在backend目录中运行此脚本"
    exit 1
fi

if [ ! -f "../web/index.html" ]; then
    echo "错误: 找不到index.html文件"
    exit 1
fi

echo "检查完成，启动服务器..."
echo "访问地址: http://localhost:8080"
echo "API地址: http://localhost:8080/api/v1"
echo "健康检查: http://localhost:8080/health"
echo ""

./crm-server