package GET

import (
	"Anzu_WebApi/ApiHandle"
	"Anzu_WebApi/Config"
	"Anzu_WebApi/Database"
	"Anzu_WebApi/Timer"
	"Anzu_WebApi/Types"
	"strconv"

	"github.com/gin-gonic/gin"
)

var cfg *Config.Config

func init() {
	cfg = Config.GetAppConfig()
}
func RootPath(c *gin.Context) {
	c.IndentedJSON(200, Types.Result{
		Code:    200,
		Message: "Hi, Takakura Anzu here!",
		Data:    "",
	})
}
func AllSubs(c *gin.Context) {
	if !ApiHandle.VerifyToken(c) {
		return
	}
	user, e := strconv.ParseInt(c.Query("user"), 10, 64)
	if e != nil {
		ApiHandle.ParseError(c, 400, "Invalid user id", e)
		return
	}
	data, err := Database.GetAllSubs(user)
	if err != nil {
		ApiHandle.ParseError(c, 500, "SQL ", err)
		return
	}
	c.IndentedJSON(200, Types.Result{
		Code:    200,
		Message: "",
		Data:    data,
	})
}
func Updates(c *gin.Context) {
	if !ApiHandle.VerifyToken(c) {
		return
	}
	Timer.UpdateRssNow()
	user, e := strconv.ParseInt(c.Query("user"), 10, 64)
	if e != nil {
		ApiHandle.ParseError(c, 400, "Invalid user id", e)
		return
	}
	data, err := Database.GetUpdates(user)
	if err != nil {
		ApiHandle.ParseError(c, 500, "SQL ", err)
		return
	}
	c.IndentedJSON(200, Types.Result{
		Code:    200,
		Message: "",
		Data:    data,
	})

}
