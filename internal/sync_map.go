package internal

import "sync"

type SyncMap[K comparable, V any] struct {
	m map[K]V
	l sync.RWMutex
}

func NewSyncMap[K comparable, V any]() *SyncMap[K, V] {
	return &SyncMap[K, V]{
		m: make(map[K]V),
		l: sync.RWMutex{},
	}
}

func (m *SyncMap[K, V]) Get(key K) V {
	m.l.RLock()
	defer m.l.RUnlock()
	return m.m[key]
}

func (m *SyncMap[K, V]) GetOk(key K) (V, bool) {
	m.l.RLock()
	defer m.l.RUnlock()
	value, ok := m.m[key]
	return value, ok
}

func (m *SyncMap[K, V]) Set(key K, value V) {
	m.l.Lock()
	defer m.l.Unlock()
	m.m[key] = value
}

func (m *SyncMap[K, _]) Delete(key K) {
	m.l.Lock()
	delete(m.m, key)
	m.l.Unlock()
}
