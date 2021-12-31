/*
# -*- coding: utf-8 -*-
# @Author : joker
# @Time : 2021/12/31 3:03 下午
# @File : pubsub.go
# @Description :
# @Attention :
*/
package listener

import "fmt"

type operation int

const (
	sub operation = iota
	subOnce
	subOnceEach
	pub
	tryPub
	unsub
	unsubAll
	closeTopic
	shutdown
)

type PubSub struct {
	cmdChan  chan cmd
	capacity int
}

type cmd struct {
	op     operation
	topics []string
	ch     chan interface{}
	msg    interface{}
}

func New(capacity int) *PubSub {
	ps := &PubSub{make(chan cmd), capacity}
	return ps
}

func (ps *PubSub) Sub(topics ...string) chan interface{} {
	return ps.sub(sub, topics...)
}

func (ps *PubSub) SubOnce(topics ...string) chan interface{} {
	return ps.sub(subOnce, topics...)
}

func (ps *PubSub) SubOnceEach(topics ...string) chan interface{} {
	return ps.sub(subOnceEach, topics...)
}

func (ps *PubSub) sub(op operation, topics ...string) chan interface{} {
	ch := make(chan interface{}, ps.capacity)
	ps.cmdChan <- cmd{op: op, topics: topics, ch: ch}
	return ch
}

func (ps *PubSub) AddSub(ch chan interface{}, topics ...string) {
	ps.cmdChan <- cmd{op: sub, topics: topics, ch: ch}
}

func (ps *PubSub) AddSubOnceEach(ch chan interface{}, topics ...string) {
	ps.cmdChan <- cmd{op: subOnceEach, topics: topics, ch: ch}
}

func (ps *PubSub) Pub(msg interface{}, topics ...string) {
	ps.cmdChan <- cmd{op: pub, topics: topics, msg: msg}
}

func (ps *PubSub) TryPub(msg interface{}, topics ...string) {
	ps.cmdChan <- cmd{op: tryPub, topics: topics, msg: msg}
}

func (ps *PubSub) Unsub(ch chan interface{}, topics ...string) {
	if len(topics) == 0 {
		ps.cmdChan <- cmd{op: unsubAll, ch: ch}
		return
	}

	ps.cmdChan <- cmd{op: unsub, topics: topics, ch: ch}
}

func (ps *PubSub) Close(topics ...string) {
	ps.cmdChan <- cmd{op: closeTopic, topics: topics}
}

func (ps *PubSub) Stop() {
	ps.cmdChan <- cmd{op: shutdown}
}

func (ps *PubSub) start() {
	reg := registry{
		topics:    make(map[string]map[chan interface{}]subType),
		revTopics: make(map[chan interface{}]map[string]bool),
	}

loop:
	for cmd := range ps.cmdChan {
		if cmd.topics == nil {
			switch cmd.op {
			case unsubAll:
				reg.removeChannel(cmd.ch)

			case shutdown:
				break loop
			}

			continue loop
		}

		for _, topic := range cmd.topics {
			switch cmd.op {
			case sub:
				reg.add(topic, cmd.ch, normal)

			case subOnce:
				reg.add(topic, cmd.ch, onceAny)

			case subOnceEach:
				reg.add(topic, cmd.ch, onceEach)

			case tryPub:
				reg.sendNoWait(topic, cmd.msg)

			case pub:
				reg.send(topic, cmd.msg)

			case unsub:
				reg.remove(topic, cmd.ch)

			case closeTopic:
				reg.removeTopic(topic)
			}
		}
	}

	for topic, chans := range reg.topics {
		for ch := range chans {
			reg.remove(topic, ch)
		}
	}
}

type registry struct {
	topics    map[string]map[chan interface{}]subType
	revTopics map[chan interface{}]map[string]bool
}

type subType int

const (
	onceAny subType = iota
	onceEach
	normal
)

func (reg *registry) add(topic string, ch chan interface{}, st subType) {
	if reg.topics[topic] == nil {
		reg.topics[topic] = make(map[chan interface{}]subType)
	}
	reg.topics[topic][ch] = st

	if reg.revTopics[ch] == nil {
		reg.revTopics[ch] = make(map[string]bool)
	}

	if reg.revTopics[ch][topic] {
		panic(fmt.Sprintf("重复:%s", topic))
	}
	reg.revTopics[ch][topic] = true
}

func (reg *registry) send(topic string, msg interface{}) {
	for ch, st := range reg.topics[topic] {
		if nil == msg {
			for topic := range reg.revTopics[ch] {
				reg.remove(topic, ch)
			}
		} else {
			ch <- msg
			switch st {
			case onceAny:
				for topic := range reg.revTopics[ch] {
					reg.remove(topic, ch)
				}
			case onceEach:
				reg.remove(topic, ch)
			}
		}
	}
}

func (reg *registry) sendNoWait(topic string, msg interface{}) {
	for ch, st := range reg.topics[topic] {
		select {
		case ch <- msg:
			switch st {
			case onceAny:
				for topic := range reg.revTopics[ch] {
					reg.remove(topic, ch)
				}
			case onceEach:
				reg.remove(topic, ch)
			}
		default:
		}

	}
}

func (reg *registry) removeTopic(topic string) {
	for ch := range reg.topics[topic] {
		reg.remove(topic, ch)
	}
}

func (reg *registry) removeChannel(ch chan interface{}) {
	for topic := range reg.revTopics[ch] {
		reg.remove(topic, ch)
	}
}

func (reg *registry) remove(topic string, ch chan interface{}) {
	if _, ok := reg.topics[topic]; !ok {
		return
	}

	if _, ok := reg.topics[topic][ch]; !ok {
		return
	}

	delete(reg.topics[topic], ch)
	delete(reg.revTopics[ch], topic)

	if len(reg.topics[topic]) == 0 {
		delete(reg.topics, topic)
	}

	if len(reg.revTopics[ch]) == 0 {
		close(ch)
		delete(reg.revTopics, ch)
	}
}
