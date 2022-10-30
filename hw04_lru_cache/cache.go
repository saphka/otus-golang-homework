package hw04lrucache

type Key string

type Cache interface {
	Set(key Key, value interface{}) bool
	Get(key Key) (interface{}, bool)
	Clear()
}

type lruCache struct {
	capacity int
	queue    List
	items    map[Key]*ListItem
}

type cacheItem struct {
	key   Key
	value interface{}
}

func (c *lruCache) Set(key Key, value interface{}) bool {
	item, ok := c.items[key]
	if ok {
		item.Value = cacheItem{
			key:   key,
			value: value,
		}
		c.queue.MoveToFront(item)
	} else {
		if c.queue.Len() == c.capacity {
			delete(c.items, c.queue.Back().Value.(cacheItem).key)
			c.queue.Remove(c.queue.Back())
		}
		c.items[key] = c.queue.PushFront(cacheItem{
			key:   key,
			value: value,
		})
	}
	return ok
}

func (c *lruCache) Get(key Key) (interface{}, bool) {
	item, ok := c.items[key]
	var val interface{}
	if ok {
		val = item.Value.(cacheItem).value
		c.queue.MoveToFront(item)
	}
	return val, ok
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
