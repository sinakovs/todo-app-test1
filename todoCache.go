package main

import "sync"

type TodoCache struct {
	mu    sync.RWMutex
	items map[string]todo
}

func NewTodoCache() *TodoCache {
	return &TodoCache{
		items: make(map[string]todo),
	}
}

func (c *TodoCache) Get(id string) (todo, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	t, ok := c.items[id]
	return t, ok
}

func (c *TodoCache) Set(t todo) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[string(t.TodosNumber)] = t
}

func (c *TodoCache) Delete(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, id)
}

func (c *TodoCache) All() []todo {
	c.mu.RLock()
	defer c.mu.RUnlock()
	todos := make([]todo, 0, len(c.items))
	for _, t := range c.items {
		todos = append(todos, t)
	}
	return todos
}
