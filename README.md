# lori 

[![Build Status](https://github.com/cr-mao/lori/workflows/Go/badge.svg)](https://github.com/cr-mao/lori/actions)
[![Go Reference](https://pkg.go.dev/badge/github.com/cr-mao/lori.svg)](https://pkg.go.dev/github.com/cr-mao/lori)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

### 1.介绍
lori 是一款基于golang的分布式web服务器框架,目标是快速构建服务。 
- http server 基于gin 
- grpc server  
- grpc client 


### 2.安装
```shell
go get github.com/cr-mao/lori@v0.0.3
```


### 3.模块化设计 

服务注册发现 
- consul 
```shell
go get github.com/cr-mao/lori/registry/consul@v0.0.3
```

指标监控
- prometheus

日志
- zap 



### 4. 参考
- kratos
- go-zero
- due框架
- iam 极客时间go语言项目实战


