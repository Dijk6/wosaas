// Package eventbus 进程内事件总线。
package eventbus

import (
	"context"
	"runtime/debug"
	"sync"

	"github.com/sirupsen/logrus"
)

// Event 所有事件必须实现该接口。
type Event interface {
	WithContext(ctx context.Context)
	GetContext() context.Context
}

// HandlerFunc 事件处理函数。
type HandlerFunc func(Event)

// EventBus 维护事件类型与订阅者列表。
type EventBus struct {
	subscribers map[string][]HandlerFunc
	lock        sync.RWMutex
	logger      *logrus.Logger
}

// New 创建事件总线。
func New() *EventBus {
	return &EventBus{
		subscribers: make(map[string][]HandlerFunc),
		logger:      logrus.StandardLogger(),
	}
}

// SetLogger 设置日志实例。
func (eb *EventBus) SetLogger(logger *logrus.Logger) {
	if logger != nil {
		eb.logger = logger
	}
}

// Register 注册事件处理器。
func (eb *EventBus) Register(eventType string, handler HandlerFunc) {
	eb.lock.Lock()
	defer eb.lock.Unlock()
	eb.subscribers[eventType] = append(eb.subscribers[eventType], handler)
}

// Publish 异步发布事件，每个 handler 单独协程执行并 recover。
func (eb *EventBus) Publish(eventType string, event Event) {
	eb.lock.RLock()
	handlers := append([]HandlerFunc(nil), eb.subscribers[eventType]...)
	eb.lock.RUnlock()

	for _, handler := range handlers {
		go func(h HandlerFunc) {
			defer func() {
				if r := recover(); r != nil {
					eb.logger.WithFields(logrus.Fields{
						"event_type": eventType,
						"panic":      r,
						"stack":      string(debug.Stack()),
					}).Error("event handler panic")
				}
			}()
			h(event)
		}(handler)
	}
}

// BaseEvent 提供 Event 接口的默认实现，业务事件可内嵌它。
type BaseEvent struct {
	ctx context.Context
}

func (e *BaseEvent) WithContext(ctx context.Context) {
	e.ctx = ctx
}

func (e *BaseEvent) GetContext() context.Context {
	if e.ctx == nil {
		return context.Background()
	}
	return e.ctx
}

// ExampleEvent 事件总线示例事件。
type ExampleEvent struct {
	BaseEvent
	Message string
}
