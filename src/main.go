package main

import (
	"common/utils"
	"data"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"http/handler"
	"http/middleware"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"version"
)

func adminGroup( userRouter *gin.RouterGroup) {

	userRouter.StaticFS(data.Config.StaticPath, http.Dir(data.Config.PATH))
	adminRouter := userRouter.Group("")
	adminRouter.Use(middleware.LoginRequire())
	{
		adminRouter.POST("poker.gm/upload.do", handler.UploadIcon)
		adminRouter.POST("upload", handler.Upload)
	}
}

func main() {
	if len(os.Args)>1{
		if os.Args[1]=="ver"{
			fmt.Println("GameServer Version:",version.VERSION)
		}else{
			fmt.Println(`
	gameServer ver  //output version
`)
		}
		return
	}

	utils.LogFlag.Add(utils.ALL_PRINT | utils.TRACE_WRITE)
	defer utils.RecoverHandle("resource server over ...")
	utils.WLog("###版本 0.0.1")
	c := make(chan os.Signal)
	signal.Notify(c)
	//监听指定信号
	signal.Notify(c, syscall.SIGHUP, os.Interrupt)
	go startServer()

	//阻塞直至有信号传入
	for {
		sig := <-c
		if strings.Contains(sig.String(), "termi") || strings.Contains(sig.String(), "interrupt") {
			utils.WLog("receive signal interrupt string wai：", sig.String())
			os.Exit(1)
		}
	}
}

func startServer() {
	defer utils.RecoverHandle("resource server start ...")

	if ok := data.Init(); !ok {
		utils.WErr("data init err !")
		return
	}

	if data.Config.RunMode != "dev" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(gin.Logger())
	if data.Config.RunMode == "dev" {
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	router.Use(middleware.Cors())

	group := router.Group("")

	adminGroup(group)

	err := router.Run(data.Config.Port)
	if err != nil {
		utils.WErr("server run err.", err.Error())
		panic(err.Error())
	}

}
