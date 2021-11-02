package persister

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type Persister struct {
	mu       sync.Mutex
	metadata map[string][]byte
	id       int64
}

func NewPersister() *Persister {
	return &Persister{
		metadata: make(map[string][]byte),
	}
}

func (p *Persister) Set(key string, value []byte) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.metadata[key] = value
}

func (p *Persister) Add(value []byte) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.metadata[strconv.FormatInt(p.id, 10)] = value
	p.id++
}

func (p *Persister) Get(key string) ([]byte, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if m, exists := p.metadata[key]; exists {
		return m, nil
	}
	return nil, errors.New(fmt.Sprintf("Key %s not found", key))
}

func (p *Persister) Delete(key string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	delete(p.metadata, key)
}

func (p *Persister) Clear() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.metadata = make(map[string][]byte)
}

func (p *Persister) GetAll() map[string][]byte {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.metadata
}
