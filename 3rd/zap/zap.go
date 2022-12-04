package main

import (
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"go-demo/common"
	"go.uber.org/zap"
)

func zapDemo() {
	common.Delimiter("for zap")
	// zap.S().Warn("fatal occur")
	// zap.S().Warnf("fatal occur")
	// zap.S().Fatalf("fatal occur")

	var val atomic.Value
	// val.Store("str")
	if val.Load() == nil {
		fmt.Println("hit")
		val.Store("str1")
	} else {
		fmt.Println(val)
		fmt.Println(val.Load().(string))
	}
	fmt.Println(val.Load().(string))

	logger := zap.NewExample()
	defer logger.Sync()

	url := "https://example.org/zap"

	logger.Info("failed msg",
		zap.String("url", url),
		zap.Int("attempt", 3),
		zap.Duration("backoff", time.Second))
	logger.Warn("fatal")
	sugar := logger.Sugar()

	sugar.Infow("failed msg",
		"url", url,
		"attempt", 3,
		"backoff", time.Second,
	)
	sugar.Infof("failed to fetch url.")
}

func main() {
	zapDemo()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		for {
			fmt.Println("aaa")
			time.Sleep(3 * time.Second)
		}
	}()
	go func() {
		time.Sleep(9 * time.Second)
		fmt.Println("GC start")
		runtime.GC()
		fmt.Println("GC end")
	}()
	wg.Wait()
}
