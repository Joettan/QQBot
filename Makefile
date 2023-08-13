# 二进制文件的名字
BINARY_NAME=my_service

# 执行编译
build:
	@echo "Building..."
	@go version
	go build -o $(BINARY_NAME) ./cmd/main.go

# 启动服务
start:
	@echo "Starting service..."
	./$(BINARY_NAME) &

# 停止服务
stop:
	@echo "Stopping service..."
	pkill $(BINARY_NAME) || true

# 清除生成的二进制文件
clean:
	@echo "Cleaning up..."
	rm -f $(BINARY_NAME)

.PHONY: build start stop clean
