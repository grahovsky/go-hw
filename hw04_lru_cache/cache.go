package hw04lrucache

import "sync"

type Key string

type Item struct {
	Key   Key
	Value interface{}
}

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	sync.Mutex
	queue List
	items map[Key]*ListItem
}

func NewCache(capacity int) Cache {
	return &lruCache{
		capacity: capacity,
		queue:    NewList(),
		items:    make(map[Key]*ListItem, capacity),
	}
}

func (lru *lruCache) Set(key Key, value interface{}) bool {
	lru.Lock()
	defer lru.Unlock()

	if item, exists := lru.items[key]; exists {
		lru.queue.MoveToFront(item)
		item.Value.(*Item).Value = value
		return true
	}

	if lru.queue.Len() >= lru.capacity {
		removeItem := lru.queue.Back()
		lru.queue.Remove(removeItem)
		delete(lru.items, removeItem.Value.(*Item).Key)
	}

	lru.items[key] = lru.queue.PushFront(&Item{key, value})
	return false
}

func (lru *lruCache) Get(key Key) (interface{}, bool) {
	lru.Lock()
	defer lru.Unlock()
	if item, exists := lru.items[key]; exists {
		lru.queue.MoveToFront(item)
		return item.Value.(*Item).Value, true
	}

	return nil, false
}

func (lru *lruCache) Clear() {
	lru.Lock()
	lru.queue = NewList()
	lru.items = make(map[Key]*ListItem, lru.capacity)
	lru.Unlock()
}
