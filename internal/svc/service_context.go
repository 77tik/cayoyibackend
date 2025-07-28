package svc

import (
	"cayoyibackend/internal/config"
	"sync"
)

type ServiceContext struct {
	// 配置相关
	Config config.Config

	// 应用相关
	mu   sync.RWMutex
	Keys map[string]any
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}

func (svc *ServiceContext) Get(key string) (value any, exists bool) {
	svc.mu.RLock()
	defer svc.mu.RUnlock()

	value, exists = svc.Keys[key]
	return
}

func (svc *ServiceContext) GetString(key string) (s string) {
	if val, ok := svc.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}
