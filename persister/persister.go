package persister

import "sync"

type Persister struct {
	mu       sync.Mutex
	metadata map[string][]byte
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

func (p *Persister) Get(key string) []byte {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.metadata[key]
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
