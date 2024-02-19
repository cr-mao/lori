#!/bin/bash

docker run --restart=always -d -p 8500:8500 -p 8300:8300 -p 8301:8301 -p 8302:8302 -p 8600:8600/udp bitnami/consul:latest consul agent -dev -client=0.0.0.0
# 验证
# dig @192.168.3.109 -p 8600 consul.service.consul SRV

