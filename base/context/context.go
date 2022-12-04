// Copyright © 2022 hops. All rights reserved.

package main

import (
	"context"
	"strings"
	"time"

	"go.uber.org/zap"
)

var logger, _ = zap.NewDevelopment()
var sugar = logger.Sugar()

func watch(ctx context.Context, name string) {
	for {
		select {
		case <-ctx.Done():
			sugar.Info(name, " cancel signal arrived when ", time.Now().String())
			return
		default:
			sugar.Info(name, " get signal when ", time.Now().String())
			time.Sleep(1 * time.Second)
		}
	}
}
func do(ctx context.Context) {
	time.Sleep(10 * time.Second)
	sugar.Info("[do] end", time.Now().String())
}
func main() {
	var aa []int
	sugar.Info(len(aa))

	ctxTimeout, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	go func() {
		do(ctxTimeout)
	}()
	select {
	case <-ctxTimeout.Done():
		sugar.Info("done arrived", time.Now().String())
	case <-time.After(3 * time.Second):
		sugar.Info("timeout arrived", time.Now().String())
		cancel()
	}

	ctxTimeout, cancel = context.WithTimeout(context.Background(), 3*time.Second)

	go watch(ctxTimeout, "watch 1")
	go watch(ctxTimeout, "watch 2")
	sugar.Info(strings.Repeat("-", 10), "start wait for 8s, time=", time.Now().Unix())
	time.Sleep(8 * time.Second)
	cancel()
	sugar.Info(strings.Repeat("-", 10), "end wait, time=", time.Now().Unix())
	ctx, cancel := context.WithCancel(context.Background())

	messages := make(chan int, 10)
	defer close(messages)

	for i := 0; i < 6; i++ {
		messages <- i
	}

	for i := 1; i <= 2; i++ {
		go sub_process(i, ctx, messages)
	}
	time.Sleep(3 * time.Second)
	cancel()
	time.Sleep(2 * time.Second)
	sugar.Info("current process exit.")

	subCtx, c := context.WithTimeout(ctxTimeout, 5*time.Second)
	select {
	case <-subCtx.Done():
		sugar.Info("subCtx done")
	case <-time.After(10 * time.Second):
		c()
		sugar.Info("time stop")
	}
}

func sub_process(index int, ctx context.Context, message <-chan int) {
	newCtx, _ := context.WithCancel(ctx)
	go sub_job(index, "a", newCtx)
	go sub_job(index, "b", newCtx)
Consume:
	for {
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			sugar.Infof("[sub_process:%d]main goroutine notify to current sub process when %v.\n", index, time.Now())
			break Consume
		default:
			sugar.Infof("[sub_process:%d]receive message: %d\n", index, <-message)
		}
	}
}

func sub_job(parent int, name string, ctx context.Context) {
	for {
		time.Sleep(1 * time.Second)
		select {
		case <-ctx.Done():
			sugar.Infof("[sub_process:%d->%s] has been done. \n", parent, name)
			return
		default:
			sugar.Infof("[sub_process:%d->%s] under doing. \n", parent, name)
		}
	}
}
