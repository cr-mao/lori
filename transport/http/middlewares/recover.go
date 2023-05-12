package middlewares

import (
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"

	"github.com/cr-mao/lori/log"
	"github.com/gin-gonic/gin"
)

// Recovery  来记录 Panic 和 call stack
func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取用户的请求信息
				httpRequest, _ := httputil.DumpRequest(c.Request, true)
				// 链接中断，客户端中断连接为正常行为，不需要记录堆栈信息
				var brokenPipe bool
				if ne, ok := err.(*net.OpError); ok {
					if se, ok := ne.Err.(*os.SyscallError); ok {
						errStr := strings.ToLower(se.Error())
						if strings.Contains(errStr, "broken pipe") || strings.Contains(errStr, "connection reset by peer") {
							brokenPipe = true
						}
					}
				}
				// 链接中断的情况
				if brokenPipe {
					log.Errorf("urlpath:%s,err:%+v,request:%s", c.Request.URL.Path, err, string(httpRequest))
					c.Error(err.(error))
					c.Abort()
					// 链接已断开，无法写状态码
					return
				}
				log.Errorf("recovery from panic,err:%+v,request:%s,stacktrace:%s", c.Request.URL.Path, err, string(httpRequest))
				// 返回 500 状态码
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"message": "服务器内部错误，请稍后再试",
				})
			}
		}()
		c.Next()
	}
}
