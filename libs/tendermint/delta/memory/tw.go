/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/3/24 10:36 上午
# @File : tw.go
# @Description :
# @Attention :
*/
package memory

import (
	"container/list"
	"time"
)

type Job func(interface{})

// FOR LOCAL DEBUG TEST
type TimeWheel struct {
	interval          time.Duration
	ticker            *time.Ticker
	slots             []*list.List
	timer             map[interface{}]int
	currentPos        int
	slotNum           int
	job               Job
	addTaskChannel    chan Task
	removeTaskChannel chan interface{}
	stopChannel       chan bool
}

type Task struct {
	delay  time.Duration
	circle int
	key    interface{}
	data   interface{}
}

func New(interval time.Duration, slotNum int, job Job) *TimeWheel {
	if interval <= 0 || slotNum <= 0 || job == nil {
		return nil
	}
	tw := &TimeWheel{
		interval:          interval,
		slots:             make([]*list.List, slotNum),
		timer:             make(map[interface{}]int),
		currentPos:        0,
		job:               job,
		slotNum:           slotNum,
		addTaskChannel:    make(chan Task),
		removeTaskChannel: make(chan interface{}),
		stopChannel:       make(chan bool),
	}

	tw.initSlots()

	return tw
}

func (tw *TimeWheel) initSlots() {
	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = list.New()
	}
}

func (tw *TimeWheel) Start() {
	tw.ticker = time.NewTicker(tw.interval)
	go tw.start()
}

func (tw *TimeWheel) Stop() {
	tw.stopChannel <- true
}

func (tw *TimeWheel) AddTimer(delay time.Duration, key interface{}, data interface{}) {
	if delay < 0 {
		return
	}
	tw.addTaskChannel <- Task{delay: delay, key: key, data: data}
}

func (tw *TimeWheel) RemoveTimer(key interface{}) {
	if key == nil {
		return
	}
	tw.removeTaskChannel <- key
}

func (tw *TimeWheel) start() {
	for {
		select {
		case <-tw.ticker.C:
			tw.tickHandler()
		case task := <-tw.addTaskChannel:
			tw.addTask(&task)
		case key := <-tw.removeTaskChannel:
			tw.removeTask(key)
		case <-tw.stopChannel:
			tw.ticker.Stop()
			return
		}
	}
}

func (tw *TimeWheel) tickHandler() {
	l := tw.slots[tw.currentPos]
	tw.scanAndRunTask(l)
	if tw.currentPos == tw.slotNum-1 {
		tw.currentPos = 0
	} else {
		tw.currentPos++
	}
}

func (tw *TimeWheel) scanAndRunTask(l *list.List) {
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.circle > 0 {
			task.circle--
			e = e.Next()
			continue
		}

		go tw.job(task.data)
		next := e.Next()
		l.Remove(e)
		if task.key != nil {
			delete(tw.timer, task.key)
		}
		e = next
	}
}

func (tw *TimeWheel) addTask(task *Task) {
	pos, circle := tw.getPositionAndCircle(task.delay)
	task.circle = circle

	tw.slots[pos].PushBack(task)

	if task.key != nil {
		tw.timer[task.key] = pos
	}
}

func (tw *TimeWheel) getPositionAndCircle(d time.Duration) (pos int, circle int) {
	delaySeconds := int(d.Seconds())
	intervalSeconds := int(tw.interval.Seconds())
	circle = int(delaySeconds / intervalSeconds / tw.slotNum)
	pos = int(tw.currentPos+delaySeconds/intervalSeconds) % tw.slotNum

	return
}

func (tw *TimeWheel) removeTask(key interface{}) {
	position, ok := tw.timer[key]
	if !ok {
		return
	}
	l := tw.slots[position]
	for e := l.Front(); e != nil; {
		task := e.Value.(*Task)
		if task.key == key {
			delete(tw.timer, task.key)
			l.Remove(e)
		}

		e = e.Next()
	}
}
