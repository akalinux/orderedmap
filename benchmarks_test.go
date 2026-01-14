package omap

import (
	"cmp"
	"testing"
)

func BenchmarkPut8192(b *testing.B) {
	for range b.N {
		m := NewSliceTree[uint64, uint64](8192, cmp.Compare)
		//m := New[uint64, uint64](cmp.Compare)
		for i := range 8192 {
			v := uint64(i)
			m.Put(v, v)
		}
		if len(m.Slices) != 8192 {
			b.Fail()
		}
	}
}

func BenchmarkMapPut8192(b *testing.B) {
	for range b.N {
		m := make(map[uint64]uint64, 8192)
		for i := range 8192 {
			v := uint64(i)
			m[v] = v
		}
		if len(m) != 8192 {
			b.Fail()
		}
	}
}

func BenchmarkGet8192(b *testing.B) {
	m := New[uint64, uint64](cmp.Compare)
	for i := range 8192 {
		v := uint64(i)
		m.Put(v, v)
	}
	for range b.N {
		//m := NewSliceTree[uint64, uint64](8192, cmp.Compare)
		for i := range 8192 {
			c := uint64(i)
			if v, ok := m.Get(c); !ok || v != uint64(c) {
				b.Fail()
			}
		}
	}
}

func BenchmarkGetMap8192(b *testing.B) {
	m := make(map[uint64]uint64)
	for i := range 8192 {
		v := uint64(i)
		m[v] = v
	}
	for range b.N {
		//m := NewSliceTree[uint64, uint64](8192, cmp.Compare)
		for i := range 8192 {
			c := uint64(i)
			if v, ok := m[c]; !ok || v != uint64(c) {
				b.Fail()
			}
		}
	}
}
