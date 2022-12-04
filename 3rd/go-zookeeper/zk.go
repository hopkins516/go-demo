// Copyright © 2022 UCloud. All rights reserved.

package main

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-zookeeper/zk"
)

func callback(event zk.Event) {
	// zk.EventNodeCreated
	// zk.EventNodeDeleted
	fmt.Println("###########################")
	fmt.Println("path: ", event.Path)
	fmt.Println("type: ", event.Type.String())
	fmt.Println("state: ", event.State.String())
	fmt.Println("---------------------------")
}

type watchType int

const (
	watchTypeData = iota
	watchTypeExist
	watchTypeChild
)

type watchPathType struct {
	path  string
	wType watchType
}

func demo() {

	zch := make(chan zk.Event, 1)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			zch <- zk.Event{}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 10; i++ {
			zch <- zk.Event{}
		}
	}()
	wg.Wait()
}

var (
	path = "/NS/demo/go-zk"
)

func main() {
	// demo()
	// ch <- 100  此处不死锁
	var _ch = make(chan int, 1)
	_ch <- 1
	close(_ch)
	_ch = nil
	_ch <- 1
	if v, ok := <-_ch; !ok {
		fmt.Println("##", "read failed")
	} else {
		fmt.Println("##", v)
	}
	if v, ok := <-_ch; !ok {
		fmt.Println("##", "read failed")
	} else {
		fmt.Println("##", v)
	}
	// eventCallbackOption := zk.WithEventCallback(callback)
	c, _, err := zk.Connect([]string{"127.0.0.1:2181"}, 30*time.Second /*eventCallbackOption*/) // *10)
	if err != nil {
		panic(err)
	}
	content, stat, childrenW, err := c.ChildrenW(path)
	cs := ""
	for _, c := range content {
		cs += c + "|"
	}
	fmt.Printf(time.Now().String(), "[[%+v]] stat:[[%+v]]\n", cs, stat)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for {
			if childrenW == nil {
				fmt.Println(time.Now().String(), "childrenW is nil")
				break
			}
			select {
			case w, ok := <-childrenW:
				if !ok {
					fmt.Println(time.Now().String(), "read childrenW failed")
				} else {
					fmt.Println(time.Now().String(), "childrenW:", w)
				}
			}

			time.Sleep(3 * time.Second)
		}
	}()
	flag := int32(zk.FlagEphemeral)
	acl := zk.WorldACL(zk.PermAll)
	mypath := path + "/uboot"
	paths := strings.Split(mypath, "/")
	var parentPath string
	for _, p := range paths[1 : len(paths)-1] {
		parentPath += "/" + p
		exist, _, err := c.Exists(parentPath)
		if err != nil {
			fmt.Println(parentPath, " parent path exist")
			continue
		}
		if !exist {
			_, err = c.Create(parentPath, nil, 0, acl)
			if err != nil {
				fmt.Println(parentPath, " create parent path error")
				continue
			}
		}
	}
	result, err := c.Create(mypath, []byte("test"), flag, acl)
	fmt.Println(time.Now().String(), "create result:", result, "error", err)

	content, stat, childrenW, err = c.ChildrenW(path)
	cs = ""
	for _, c := range content {
		cs += c + "|"
	}
	fmt.Printf(time.Now().String(), "[[%+v]] stat:[[%+v]]\n", cs, stat)

	_, stat, znodeW, err := c.ExistsW(path)
	go func() {
		defer wg.Done()
		for {
			if znodeW == nil {
				fmt.Println(time.Now().String(), "znodeW is nil")
				break
			}
			select {
			case w, ok := <-znodeW:
				if !ok {
					fmt.Println(time.Now().String(), "read znodeW failed")
				} else {
					fmt.Println(time.Now().String(), "znodeW:", w)
				}
			}

			time.Sleep(3 * time.Second)
		}
	}()
	wg.Wait()
}
