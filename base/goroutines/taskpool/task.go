// Copyright © 2022 UCloud. All rights reserved.

package main

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"google.golang.org/protobuf/proto"
)

type BaseRunnable struct {
}

func (r *BaseRunnable) Run() {
	fmt.Println("base running")
}

type DerivedRunnable struct {
	BaseRunnable
	Name string
}

func (r *DerivedRunnable) Run() {
	fmt.Println(r.Name, "running")
}

type TaskHandler func(msg proto.Message)

type Task interface {
	Run()
	Timeout() chan struct{}
	GetTaskId() string
}

type ProbeUpmState struct {
	taskId   string // message session? 重复?
	handler  TaskHandler
	timeout  int64
	interval int64
	Msg      proto.Message
}

func (t *ProbeUpmState) Run() {
	t.handler(t.Msg)
}

func (t *ProbeUpmState) GetTaskId() string {
	return t.taskId
}

func (t *ProbeUpmState) GetTimeOut() int64 {
	return t.timeout
}

func (t *ProbeUpmState) Timeout() chan struct{} {
	var closeCh = make(chan struct{})
	for {
		select {
		case <-time.After(time.Duration(t.timeout) * time.Millisecond):
			closeCh <- struct{}{}
		}
	}
}

func NewProbeUpmStateHandler() *ProbeUpmState {
	t := &ProbeUpmState{
		taskId:  uuid.NewString(),
		timeout: 100,
		Msg:     nil,
	}
	t.handler = func(msg proto.Message) {}
	return t
}

func SendTaskToQueue(t Task) {
	taskQueue <- t
}
