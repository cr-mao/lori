package middlewares

import (
	"github.com/gin-gonic/gin"
)

var Middlewares = map[string]gin.HandlerFunc{
	"recovery": Recovery(), //
	"cors":     Cors(),
	//"logger":   gin.Logger(), // gin的logger ,还是一定要让用户自己外面传才好。
}
