// Copyright © 2022 UCloud. All rights reserved.

package main

import (
	"sync"
	"time"
)

func mutexPool(wg *sync.WaitGroup) {
	var d = &DerivedRunnable{
		Name: "demo",
	}
	d.Run()

	pool := GetRoutinePool()
	// var wg sync.WaitGroup
	closeCh := make(chan struct{}, 1)
	wg.Add(3)
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		closeTicker := make(chan struct{}, 1)
		defer func() {
			wg.Done()
		}()

		go func() {
			defer func() {
				wg.Done()
			}()
			select {
			case <-time.After(30 * time.Second):
				closeTicker <- struct{}{}
				// fmt.Println("producer stop, close pool")
				return
			}
		}()
		go func() {
			defer func() {
				wg.Done()
			}()
			for {
				select {
				case <-closeTicker:

					ticker.Stop()
					// fmt.Println("producer stop, close pool")
					return
				case <-ticker.C:
					t2 := pool.PeekTask()
					if t2 != nil {
						t2.Run()
						for {
							select {
							case <-t2.Timeout():

							}
						}
					}
				}
			}
		}()
	}()

	wg.Add(1)
	go func() {
		defer func() {
			closeCh <- struct{}{}
			wg.Done()
		}()
		for i := 0; i < 4; i++ {
			t1 := NewProbeUpmStateHandler()
			pool.AddTask(t1)
			time.Sleep(5 * time.Second)
		}
	}()
	// wg.Wait()
}

var taskQueue = make(chan Task, 2000)

func bufferChannelPool() {
	for task := range taskQueue {
		task.Run()
	}
}

func main() {
	var wg sync.WaitGroup
	// wg.Add(1)
	// go func() {
	// 	defer wg.Done()
	// 	mutexPool(&wg)
	// }()

	wg.Add(1)
	go func() {
		defer wg.Done()
		bufferChannelPool()
	}()

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		for i := 0; i < 100; i++ {
			t1 := NewProbeUpmStateHandler()
			SendTaskToQueue(t1)
			// time.Sleep(5 * time.Second)
		}
	}()
	time.Sleep(6 * time.Second)

	wg.Add(1)
	go func() {
		defer func() {
			wg.Done()
		}()
		for i := 0; i < 3; i++ {
			t1 := NewProbeUpmStateHandler()
			// SendTaskToQueue(t1)
			taskQueue <- t1
			// time.Sleep(5 * time.Second)
		}
	}()

	wg.Wait()
}
