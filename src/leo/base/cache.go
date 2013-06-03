/* this is data cache
*/

package base

import (
	"sync"
	"strconv"
)

//goroutine safe
type Cache struct {
	data map[string] *Record
	lock sync.RWMutex
}

func NewCache() (cache *Cache, err error) {
	cache = new(Cache)
	err = cache.init()
	return
}

func (cache *Cache) init() error {
	cache.data = make(map[string] *Record)
	return nil
}

func (cache *Cache) Start() error {
	return nil
}

func (cache *Cache) Close() error {
	return nil
}

func (cache *Cache) Get(table string, key int) (*Record) {
	k := cache.map_key(table, key)
	if k == "" {
		return nil
	}

	cache.lock.RLock()
	v, ok := cache.data[k]
	cache.lock.RUnlock()

	if !ok {
		return nil
	}
	return v
}

func (cache *Cache) Set(table string, key int, record *Record) {
	k := cache.map_key(table, key)
	if k == "" {
		return
	}

	cache.lock.Lock()
	cache.data[k] = record
	cache.lock.Unlock()
}

func (cache *Cache) Add(table string, key int, record *Record) {
	k := cache.map_key(table, key)
	if k == "" {
		return
	}

	cache.lock.Lock()
	cache.data[k] = record
	cache.lock.Unlock()
}

func (cache *Cache) Del(table string, key int) {
	k := cache.map_key(table, key)
	if k == "" {
		return
	}

	cache.lock.Lock()
	delete(cache.data, k)
	cache.lock.Unlock()
}

func (cache *Cache) map_key(table string, key int) string {
	if table == "" || key <= 0 {
		return ""
	}
	return table + strconv.Itoa(key)
}