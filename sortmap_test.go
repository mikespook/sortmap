package sortmap

import (
	"sort"
	"testing"
)

func TestData(t *testing.T) {
	sm := New(nil)
	sm.Set(1, "a")
	sm.Set(2, "b")
	a := sm.Get(1)
	if _, ok := a.(string); !ok {
		t.Errorf("Wrong type: %T", a)
	}
	if sm.Len() != 2 {
		t.Errorf("Wrong length: %d", sm.Len())
	}
	b := sm.Delete(2)
	if _, ok := b.(string); !ok {
		t.Errorf("Wrong type: %T", b)
	}
	b = sm.Delete(2)
	if b != nil {
		t.Errorf("Wrong type: %T", b)
	}
	b = sm.Get(2)
	if b != nil {
		t.Errorf("Wrong type: %T", b)
	}
	if sm.Len() != 1 {
		t.Errorf("Wrong length: %d", sm.Len())
	}
}

func TestBy(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error("Need a panic")
		}
	}()
	sm := New(nil)
	sm.Set(1, "a")
	sm.Set(2, "b")
	sm.Less(0, 1)
}

func TestSort(t *testing.T) {
	sm := New(func(ka, va, kb, kv interface{}) bool {
		return ka.(int) < kb.(int)
	})
	if c, ok := sm.End(); ok || c != nil {
		t.Errorf("Wrong data: %v", c)
	}
	if a, ok := sm.Begin(); ok || a != nil {
		t.Errorf("Wrong data: %v", a)
	}
	sm.Set(2, "b")
	sm.Set(1, "a")
	sm.Set(3, "c")
	sort.Sort(sm)
	if c, ok := sm.End(); ok {
		if str, ok := c.(int); !ok || str != 3 {
			t.Errorf("Wrong data: %v", c)
		}
	} else {
		t.Error("Wrong data")
	}
	if n, ok := sm.Next(); ok || n != nil {
		t.Errorf("Wrong data: %v", n)
	}
	if a, ok := sm.Begin(); ok {
		if str, ok := a.(int); !ok || str != 1 {
			t.Errorf("Wrong data: %v", a)
		}
	} else {
		t.Error("Wrong data")
	}
	if n, ok := sm.Next(); !ok || n != 2 {
		t.Errorf("Wrong data: %v", n)
	}
}
