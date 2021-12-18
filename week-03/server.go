package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "golang.org/x/sync/errgroup"
)

// 题目：基于 errgroup 实现一个 http server 的启动和关闭 ，
// 以及 linux signal 信号的注册和处理，要保证能够一个退出，全部注销退出。
func main() {
    ctx := context.Background()
    g, ctx := errgroup.WithContext(ctx)

    // 启动http
    mux := &http.ServeMux{}
    mux.HandleFunc("/", hello)
    s := &http.Server{
        Addr:           ":8080",
        Handler:        mux,
        ReadTimeout:    30 * time.Second,
        WriteTimeout:   30 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    g.Go(func() error {
        if err := s.ListenAndServe(); err != nil {
            return err
        }
        return nil
    })

    g.Go(func() error { // 另外一个 goroutine，用来验证是否会关闭 errGroup的其他goroutine
        for  {
            select {
            case <-ctx.Done():
                return nil
            default:
                fmt.Print("do my job\n")
                time.Sleep(time.Second * 1)
            }
        }
    })

    go func() {
        g.Wait()
    }()

    // 接收并处理系统信号
    handleSig(s, ctx)
}

func handleSig(srv *http.Server, gCtx context.Context) {
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit // 等待接收系统信号

    // 预留服务10秒钟退出时间，超时则强制退出
    timeOutCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()
    if err := srv.Shutdown(timeOutCtx); err != nil {
        log.Printf("app force to shut down, err:%v", err)
    } else {
        log.Println("app gracefully  shut down")
    }
    gCtx.Done() // 主服务结束，其他的 goroutine 也结束
}

func hello(w http.ResponseWriter, req *http.Request) {
    time.Sleep(time.Second * 5)
    fmt.Fprintf(w, "hello\n")
}
