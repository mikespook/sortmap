package sortmap_test

import (
	"fmt"
	"github.com/mikespook/sortmap"
	"sort"
	"sync"
)

func ExampleSortMap() {
	var wg sync.WaitGroup
	by := func(ka, va, kb, vb interface{}) bool {
		return ka.(int) < kb.(int)
	}
	sm := sortmap.New(by)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			for i := 0; i < 10000; i++ {
				sm.Set(i, i)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	sort.Sort(sm)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			for a, ok := sm.Begin(); ok; a, ok = sm.Next() {
				fmt.Println(a, ok)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
