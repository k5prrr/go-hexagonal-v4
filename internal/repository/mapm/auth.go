// Package mapm provides a thread-safe map for storing full user data.
package mapm

import (
	"app/internal/app/core/domain"

	"sync"
)

type IMapAuth interface {
	Set(token string, userFull *domain.UserFull)
	Get(token string) *domain.UserFull
	Delete(token string)
}

type MapAuth struct {
	mu     sync.RWMutex
	m      map[string]*domain.UserFull
	maxLen int64
}

func NewMapAuth(maxLen int64) *MapAuth {
	return &MapAuth{
		m:      make(map[string]*domain.UserFull),
		maxLen: maxLen,
	}
}

func (m *MapAuth) Set(token string, userFull *domain.UserFull) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.maxLen > 0 && int64(len(m.m)) >= m.maxLen {
		clear(m.m)
	}

	m.m[token] = userFull
}

func (m *MapAuth) Get(token string) *domain.UserFull {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return m.m[token]
}

func (m *MapAuth) Delete(token string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.m, token)
}

