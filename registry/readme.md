后续搞成模块化的， consul 目录下搞个go.mod出来就行。 
```text
module github.com/cr-mao/lori/registry/consul

go 1.18

replace github.com/cr-mao/lori => ../../

require (
	github.com/cr-mao/lori v0.0.0-00010101000000-000000000000
	github.com/hashicorp/consul/api v1.20.0
)

```
