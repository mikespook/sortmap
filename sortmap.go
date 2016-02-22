package sortmap

import (
	"sync"
)

// CompareFunc has 4 paramaters in two pair: `a(ka=>va)` and `b(kb=>vb)`.
type CompareFunc func(ka, va, kb, vb interface{}) bool

// SortMap illustrates how to sort a map
type SortMap struct {
	// an user-defined comparation logical
	by CompareFunc

	// internal map, keeps data
	m map[interface{}]interface{}
	// internal slice, keeps order
	s []interface{}
	// position indicator
	p int
	// mutex of position indicator
	pmutex sync.Mutex

	// mutex of data
	sync.RWMutex
}

// New returns an instance of SortMap
func New(by CompareFunc) *SortMap {
	return &SortMap{
		by: by,
		m:  make(map[interface{}]interface{}),
		s:  make([]interface{}, 0),
	}
}

// Set asserts data `v` with key `k`
func (sm *SortMap) Set(k interface{}, v interface{}) {
	sm.Lock()
	_, ok := sm.m[k]
	sm.m[k] = v
	// don't add the same key repeatly
	if !ok {
		sm.s = append(sm.s, k)
	}
	sm.Unlock()
}

// Delete removes the data with key `k` and returns it
func (sm *SortMap) Delete(k interface{}) interface{} {
	sm.Lock()
	defer sm.Unlock()
	a, ok := sm.m[k]
	if ok {
		delete(sm.m, k)
		// delete from the slice
		for i, v := range sm.s {
			if v == k {
				sm.s = append(sm.s[:i], sm.s[i+1:]...)
				break
			}
		}
		return a
	}
	return nil
}

// Get returns the data with key `k`
func (sm *SortMap) Get(k interface{}) interface{} {
	sm.RLock()
	defer sm.RUnlock()
	return sm.m[k]
}

// Len is part of sort.Interface and return the lengh of data
func (sm *SortMap) Len() int {
	sm.RLock()
	defer sm.RUnlock()
	return len(sm.m)
}

// Swap is part of sort.Interface and should not be used directly
func (sm *SortMap) Swap(i, j int) {
	sm.Lock()
	defer sm.Unlock()
	sm.s[i], sm.s[j] = sm.s[j], sm.s[i]
}

// Next returns the next data and true, or nil and false if it's the end of data
func (sm *SortMap) Next() (v interface{}, has bool) {
	sm.RLock()
	defer sm.RUnlock()
	// lock for the position indicator
	sm.pmutex.Lock()
	defer sm.pmutex.Unlock()
	if sm.p+1 == len(sm.s) {
		return nil, false
	}
	sm.p++
	a := sm.s[sm.p]
	return a, true
}

// Begin resets the position indicator and return the first data and true,
// or nil and false if not data has been added into the map
func (sm *SortMap) Begin() (interface{}, bool) {
	sm.RLock()
	defer sm.RUnlock()
	if len(sm.s) == 0 {
		return nil, false
	}
	sm.pmutex.Lock()
	defer sm.pmutex.Unlock()
	sm.p = 0
	return sm.s[sm.p], true
}

// End sets the position indicator to the end of data and return it and true,
// or nil and false if not data has been added into the map
func (sm *SortMap) End() (interface{}, bool) {
	sm.RLock()
	defer sm.RUnlock()
	if len(sm.s) == 0 {
		return nil, false
	}
	sm.pmutex.Lock()
	defer sm.pmutex.Unlock()
	sm.p = len(sm.s) - 1
	return sm.s[sm.p], true
}

// Less is part of sort.Interface. It is implemented by calling the "By" closure in the sorter.
func (sm *SortMap) Less(i, j int) bool {
	sm.Lock()
	defer sm.Unlock()
	if sm.by == nil {
		panic("SortMap.By not set")
	}
	return sm.by(sm.s[i], sm.m[sm.s[i]], sm.s[j], sm.m[sm.s[j]])
}
