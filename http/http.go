package http

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echolog "github.com/onrik/logrus/echo"
	"github.com/qwqcode/qwquiver/bindata"
	"github.com/qwqcode/qwquiver/config"
	"github.com/sirupsen/logrus"
)

var api *echo.Group
var Injections = [](func(api *echo.Group)){}

// Run 运行 http server
func Run() {
	e := echo.New()
	e.HideBanner = true

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://localhost:3000"}, // For dev
	}))
	e.Logger = echolog.NewLogger(logrus.StandardLogger(), "")
	e.Use(echolog.Middleware(echolog.DefaultConfig))

	fileServer := http.FileServer(bindata.AssetFile())
	e.GET("/*", echo.WrapHandler(fileServer))

	api := e.Group("/api")
	api.GET("/query", queryHandler)
	api.GET("/query/avg", queryAvgHandler)
	api.GET("/conf", confHandler)
	api.GET("/analyze", analyzeHandler)
	api.GET("/school/all", schoolAllHandler)

	// 功能注入
	for _, inject := range config.HTTPInjections {
		inject(e, api)
	}

	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", config.Instance.Port)))
}
