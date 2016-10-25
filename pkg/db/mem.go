package db

import (
	"errors"
	"github.com/llitfkitfk/GoHighPerformance/pkg/model"
	"sync"
)

var ErrNotFound = errors.New("not found")

type Mem struct {
	mx sync.RWMutex
	m  map[string]model.Model
}

func NewMem() *Mem {
	return &Mem{m: make(map[string]model.Model)}
}

func (m *Mem) Save(key model.Key, model model.Model) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	m.m[key.String()] = model
	return nil
}

func (m *Mem) Delete(key model.Key) error {
	m.mx.Lock()
	defer m.mx.Unlock()
	delete(m.m, key.String())
	return nil
}

func (m *Mem) Get(key model.Key, model model.Model) error {
	m.mx.RLock()
	defer m.mx.RUnlock()
	md, ok := m.m[key.String()]
	if !ok {
		return ErrNotFound
	}
	return model.Set(md)
}
