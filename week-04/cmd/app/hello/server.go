package hello

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	"week-04/app/hello/router"

	"github.com/gin-gonic/gin"
)

func InitServer() *http.Server {
	r := setupRouter()
	srv := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", "127.0.0.1", 7000),
		Handler: r,
	}
	return srv
}

//运行服务
func Run() error {
	srv := InitServer()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Hello server Shutdown:", err)
	}
	log.Println("Hello server exiting")
	return nil
}

// 初始化路由
func setupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(gin.Recovery())
	r.GET("/ping", router.Ping)
	return r
}
