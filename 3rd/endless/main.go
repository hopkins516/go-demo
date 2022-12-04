// Copyright © 2022 UCloud. All rights reserved.

package main

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/fvbock/endless"
	"github.com/gin-gonic/gin"
)

func main() {
	go func() {
		for {
			log.Printf("gc start")
			runtime.GC()
			log.Printf("gc done")
			time.Sleep(5 * time.Second)
		}
	}()
	// gin
	router := gin.Default()
	router.GET("/hello", func(c *gin.Context) {
		delay := c.Query("delay")
		dt, _ := strconv.Atoi(delay)
		log.Printf("--- sleep time: %d", dt)
		time.Sleep(time.Second * time.Duration(dt))
		c.String(http.StatusOK, "world 111") // 返回值
	})

	// endless
	endless.DefaultReadTimeOut = 10 * time.Second
	endless.DefaultWriteTimeOut = 30 * time.Second // 写 超时时间为 30s
	endless.DefaultMaxHeaderBytes = 1 << 20        // 请求头最大为 1m
	endPoint := fmt.Sprintf(":%d", 8811)           // 端口

	srv := endless.NewServer(endPoint, router)
	srv.BeforeBegin = func(add string) {
		log.Printf("Actual pid is %d", syscall.Getpid())
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Printf("Server err: %v", err)
	}
}
