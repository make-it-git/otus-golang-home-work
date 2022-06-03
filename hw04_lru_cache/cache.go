package hw04lrucache

import (
	"sync"
)

type Key string

type cacheItem struct {
	key   Key
	value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
	lock     sync.RWMutex
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	c.lock.Lock()
	defer c.lock.Unlock()

	existingItem, ok := c.items[key]
	if !ok {
		item := c.queue.PushFront(cacheItem{
			key:   key,
			value: value,
		})
		c.items[key] = item
	}

	if ok {
		existingItem.Value = cacheItem{
			key:   key,
			value: value,
		}
		c.queue.MoveToFront(existingItem)
	}

	if c.queue.Len() > c.capacity {
		back := c.queue.Back()
		c.queue.Remove(back)
		deleteKey := back.Value.(cacheItem)
		delete(c.items, deleteKey.key)
	}

	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	item, ok := c.items[key]

	if ok {
		return item.Value.(cacheItem).value, ok
	}

	return nil, ok
}

func (c *lruCache) Clear() {
	c.queue = NewList()
	c.items = make(map[Key]*ListItem, c.capacity)
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}
