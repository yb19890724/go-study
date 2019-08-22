package eventhandlers

import (
	"sync"
)

// 事件名称
type Event string

// 处理接口
type Handler interface {
	Handle()
}

// 事件处理寄存器
type EventHandlers struct {
	ehsMap map[Event][]Handler
	lock   sync.RWMutex
}

// 单例方式
var inst *EventHandlers

var once sync.Once

func GetInstance() *EventHandlers {
	once.Do(func() {
		inst = &EventHandlers{ehsMap: make(map[Event][]Handler)}
	})
	return inst
}

// 发布
func (e *EventHandlers) Pub(event Event) {
	e.lock.RLock()
	defer e.lock.RUnlock()
	if handlers, ok := e.ehsMap[event]; ok {
		for _, handler := range handlers {
			go handler.Handle()
		}
	}

}

// 订阅
func (e *EventHandlers) Sub(event Event, handler Handler) {
	e.lock.Lock()
	defer e.lock.Unlock()
	e.ehsMap[event] = append(e.ehsMap[event], handler)
}
