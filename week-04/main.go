package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
	"week-04/cmd/app/admin"
	"week-04/cmd/app/hello"

	"golang.org/x/sync/errgroup"
)

var Servers = make([]*http.Server, 0)

func init() {
	Servers = append(Servers, admin.InitServer())
	Servers = append(Servers, hello.InitServer())
}

// 启动所有程序
func main() {
	ctx := context.Background()
	g, ctx := errgroup.WithContext(ctx)

	for _, s := range Servers {
		g.Go(func() error {
			if err := s.ListenAndServe(); err != nil {
				log.Print(err)
				log.Print("00000")
				return err
			}
			return nil
		})
	}

	go func() {
		err := g.Wait()
		if err != nil {
			log.Print(err)
			return
		}
	}()

	// 接收并处理系统信号
	handleSig(Servers, ctx)
}

func handleSig(servers []*http.Server, gCtx context.Context) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit // 等待接收系统信号

	// 预留服务10秒钟退出时间，超时则强制退出
	timeOutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	for _, srv := range servers {
		if err := srv.Shutdown(timeOutCtx); err != nil {
			log.Printf("app force to shut down, err:%v", err)
		} else {
			log.Println("app gracefully  shut down")
		}
	}
	gCtx.Done() // 主服务结束，其他的 goroutine 也结束
}
