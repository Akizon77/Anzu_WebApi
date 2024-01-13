package ApiHandle

import (
	"Anzu_WebApi/Config"
	"Anzu_WebApi/Types"
	"fmt"
	"github.com/gin-gonic/gin"
)

var cfg *Config.Config

func init() {
	cfg = Config.GetAppConfig()
}
func ParseError(c *gin.Context, code int, message string, err error) {
	c.IndentedJSON(code, Types.Result{
		Code:    code,
		Message: fmt.Sprint(message, ", ", err),
		Data:    nil,
	})
}
func VerifyToken(c *gin.Context) bool {
	token := c.Query("token")
	if token == cfg.Token {
		return true
	} else {
		c.IndentedJSON(401, gin.H{
			"code":    401,
			"message": "401 Unauthorized",
		})
		return false
	}
}
func NoRouter(c *gin.Context) {
	c.IndentedJSON(404, gin.H{
		"code":    404,
		"message": "404 Not Found",
	})
}
func MethodNotAllowed(c *gin.Context) {
	c.IndentedJSON(405, gin.H{
		"code":    405,
		"message": "405 Method Not Allowed",
	})
}
