VERSION = `git rev-parse --short HEAD`
BUILDTIME = `date +%FT%T`

%:
	@true


## protoc-gen-go: 安装pb文件生成工具
.PHONY: protoc-gen-go
protoc-gen-go:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

## protoc-gen-go-grpc: 安装grpc pb文件生成工具
.PHONY: protoc-gen-go-grpc
protoc-gen-go-grpc:
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2

## tidy: 整理现有的依赖
.PHONY: tidy
tidy:
	go mod tidy

## download: go依赖下载
.PHONY: download
download:
	go mod download

## vet: 静态检测全部go代码
vet:
	go vet ./...

## test: 单元测试全部测试代码
test:
	go test ./...

## help: Show this help info.
.PHONY: help
help: Makefile
	@echo  "\nUsage: make <TARGETS> <OPTIONS> ...\n\nTargets:"
	@sed -n 's/^##//p' $< | column -t -s ':' | sed -e 's/^/ /'



