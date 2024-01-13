package POST

import (
	"Anzu_WebApi/ApiHandle"
	"Anzu_WebApi/Config"
	"Anzu_WebApi/Database"
	"Anzu_WebApi/Types"
	"errors"
	"github.com/gin-gonic/gin"
)

var cfg *Config.Config

func init() {
	cfg = Config.GetAppConfig()
}
func AddNewSubs(c *gin.Context) {
	if !ApiHandle.VerifyToken(c) {
		return
	}
	req := Types.POST_Add{}
	//从Body读取数据 解析到req
	err := c.BindJSON(&req)
	//处理解析Body错误的情况
	if err != nil {
		ApiHandle.ParseError(c, 400, "Bad request", err)
		return
	}
	//处理传入数据为空的情况
	if req.User == 0 || req.Title == "" || req.Link == "" {
		ApiHandle.ParseError(c, 400, "Bad request", errors.New("user, title and link can not be null"))
		return
	}
	//数据库操作
	err = Database.AddNew(req.User, req.Title, req.Link)
	//处理数据库错误
	if err != nil {
		ApiHandle.ParseError(c, 500, "SQL Error:", err)
		return
	}
	//所有操作成功完成，返回success
	c.IndentedJSON(200, Types.Result{
		Code:    200,
		Message: "success",
		Data:    nil,
	})
}
func DelSubs(c *gin.Context) {
	if !ApiHandle.VerifyToken(c) {
		return
	}
	//与Add共用一个数据结构
	req := Types.POST_Add{}
	//从Body读取数据 解析到req
	err := c.BindJSON(&req)
	//处理解析Body错误的情况
	if err != nil {
		ApiHandle.ParseError(c, 400, "Bad request", err)
		return
	}
	//处理传入数据为空的情况
	if req.User == 0 || req.Link == "" {
		ApiHandle.ParseError(c, 400, "Bad request", errors.New("user and link can not be null"))
		return
	}
	//数据库操作
	err = Database.Del(req.User, req.Link)
	//处理数据库错误
	if err != nil {
		ApiHandle.ParseError(c, 500, "SQL Error:", err)
		return
	}
	//所有操作成功完成，返回success
	c.IndentedJSON(200, Types.Result{
		Code:    200,
		Message: "success",
		Data:    nil,
	})
}
