/* this is data cache
*/

package db

import (
	"leo/base"
)

type RecordMap map[int] *base.Record
type Cache map[string] RecordMap

func NewCache() (cache Cache, err error) {
	cache = make(Cache)
	err = cache.init()
	return
}

func (cache Cache) init() error {
	return nil
}

func (cache Cache) Start() {
}

func (cache Cache) Close() {
}

func (cache Cache) Get(table string, key int) (*base.Record) {
	if table == "" {
		return nil
	}

	mp, ok := cache[table]
	if !ok {
		return nil
	}

	rcd, ok := mp[key]
	if !ok {
		return nil
	}

	return rcd
}

func (cache Cache) Set(table string, key int, record *base.Record) {
	if table == "" {
		return
	}

	mp, ok := cache[table]
	if !ok {
		mp = make(RecordMap)
		cache[table] = mp
	}

	mp[key] = record
}

func (cache Cache) Add(table string, key int, record *base.Record) {
	if table == "" {
		return
	}

	mp, ok := cache[table]
	if !ok {
		mp = make(RecordMap)
		cache[table] = mp
	}

	mp[key] = record
}

func (cache Cache) Del(table string, key int) {
	if table == "" {
		return
	}

	mp, ok := cache[table]
	if !ok {
		return
	}

	delete(mp, key)

	if len(mp) == 0 {
		delete(cache, table)
	}
}