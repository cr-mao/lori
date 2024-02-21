package middlewares

import (
	"github.com/gin-gonic/gin"
)

// 外部可以用
const (
	RECOVERY string = "recovery"
	CORS     string = "cors"
)

var Middlewares = map[string]gin.HandlerFunc{
	"recovery": Recovery(), //
	"cors":     Cors(),
	//"logger":   gin.Logger(), // gin的logger ,还是一定要让用户自己外面传才好。
}
