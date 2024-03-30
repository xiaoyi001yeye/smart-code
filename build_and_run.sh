#!/bin/bash

# 定义 Docker 镜像名称和容器名称
IMAGE_NAME="smartcodeql"
CONTAINER_NAME="smartcodeql_container"
CONTAINER_ID_FILE="/tmp/smartcodeql_container_id" # 用于存储容器ID的文件

# 构建 Docker 镜像
docker build --pull --rm -f "Dockerfile" -t $IMAGE_NAME:latest "."

if [ $? -ne 0 ]; then
    echo "Docker build failed with exit status $?"
    exit 1
fi

# 获取当前正在运行的容器ID
running_container_id=$(docker ps -q -f name=^/$CONTAINER_NAME$)

# 如果容器正在运行，则停止并删除它
if [ -n "$running_container_id" ]; then
    echo "Stopping and removing existing container: $running_container_id"
    docker stop $running_container_id
    docker rm $running_container_id
fi

# 检查容器ID文件是否存在，如果存在则删除旧容器
if [ -f "$CONTAINER_ID_FILE" ]; then
    old_container_id=$(cat $CONTAINER_ID_FILE)
    echo "Removing old container: $old_container_id"
    docker rm $old_container_id
fi
# 运行新容器
echo "Running new container with name: $CONTAINER_NAME"
docker run -d --name $CONTAINER_NAME -p 3001:3000 $IMAGE_NAME:latest

# 将新容器的ID写入文件
echo $(docker ps -l -q) > $CONTAINER_ID_FILE

echo "Checking service status..."
response=$(curl -s -o /dev/null -w '%{http_code}' http://localhost:3001)
if [ "$response" == "200" ]; then
    echo "Service is running."
else
    echo "Service is not running. Response code: $response"
    exit 1
fi