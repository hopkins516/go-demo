// Copyright © 2022 UCloud. All rights reserved.

package main

import "sync"

type RoutinePool struct {
	Mutex       sync.Mutex
	Task        []Task
	MaxRoutines uint64
}

var once sync.Once
var _pool *RoutinePool

func GetRoutinePool() *RoutinePool {
	once.Do(func() {
		_pool = new(RoutinePool)
	})
	return _pool
}

func (p *RoutinePool) AddTask(task Task) {
	p.Mutex.Lock()
	p.Task = append(p.Task, task)
	p.Mutex.Unlock()
}

func (p *RoutinePool) PeekTask() Task {
	if len(p.Task) == 0 {
		return nil
	}
	p.Mutex.Lock()
	task := p.Task[0]
	p.Task = p.Task[1:]
	p.Mutex.Unlock()
	return task
}
