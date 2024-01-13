package main

import (
	"Anzu_WebApi/ApiHandle"
	"Anzu_WebApi/ApiHandle/GET"
	"Anzu_WebApi/ApiHandle/POST"
	"Anzu_WebApi/Config"
	"Anzu_WebApi/Log"
	"Anzu_WebApi/Timer"
	"github.com/gin-gonic/gin"
	"io"
)

var cfg *Config.Config

func init() {
	cfg = Config.GetAppConfig()
}
func main() {
	gin.SetMode(gin.ReleaseMode)
	logwww := Log.GetWebLogger()
	Log.Info("Running http server on ", cfg.Listen)
	go Timer.UpdateRssNow()
	go Timer.StartRssAutoRefresh(600)
	gin.DefaultWriter = io.Discard
	router := gin.Default()
	//日志处理事件
	router.Use(func(c *gin.Context) {
		c.Next()
		logwww.Info(c.Writer.Status(), " ", c.ClientIP(), " ", c.Request.Method, " ", c.Request.RequestURI)
	})

	router.HandleMethodNotAllowed = true
	router.NoRoute(ApiHandle.NoRouter)
	router.NoMethod(ApiHandle.MethodNotAllowed)

	router.GET("/", GET.RootPath)
	router.GET("/all", GET.AllSubs)
	router.GET("/updates", GET.Updates)
	router.POST("/add", POST.AddNewSubs)
	router.POST("/del", POST.DelSubs)

	err := router.Run(cfg.Listen)
	if err != nil {
		Log.Error("Can not run http server", err)
	}
}
