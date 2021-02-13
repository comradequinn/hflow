package store

import (
	"runtime"
	"strconv"
	"testing"
)

func TestAllocs(t *testing.T) {
	runtime.GC()

	heapSize := func() uint64 {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		return m.Alloc
	}

	Init(100)

	initHeapSize := heapSize()

	for i := 0; i < 100; i++ { // add one more than the capacity to ensure we test the reassignment of the tail
		Insert(Record{})
	}

	runtime.GC()

	capacityReachedHeapSize := heapSize()

	allocsToCapacity := capacityReachedHeapSize - initHeapSize

	if capacityReachedHeapSize <= initHeapSize {
		t.Fatalf("expected allocations during lilst expansion. init: %v, capacity-reached: %v", initHeapSize, capacityReachedHeapSize)
	}

	for i := 0; i < 1000; i++ { // add one more than the capacity to ensure we test the reassignment of the tail
		Insert(Record{})
	}

	runtime.GC()
	capacityExceededHeapSize := heapSize()

	// difficult to account for small allocations, so we check that adding 10x that of the capacity has resulted in less allocs than those made by reaching capacity;
	// demonstrating that we have considerably reduced post-capacity allocs
	if capacityExceededHeapSize > (capacityReachedHeapSize + allocsToCapacity) {
		t.Fatalf("expected no further allocations during list memory re-use phase. capacity-reached: %v,capacity-exceded: %v", capacityReachedHeapSize, capacityExceededHeapSize)
	}
}

func TestGet(t *testing.T) {
	Init(4)

	for i := 0; i < 5; i++ { // add one more than the capacity to ensure we test the reassignment of the tail and both seekdirections
		Insert(Record{Data: strconv.Itoa(i + 1)})
	}

	r, ok := Get(1000)

	if ok {
		t.Fatalf("expected record with id %v to not exist as it was never added. got %v, %v", 1000, r, ok)
	}

	r, ok = Get(1)

	if ok {
		t.Fatalf("expected record with id %v to not exist as the capacity restrictions will have cycled it out. got %v, %v", 1, r, ok)
	}

	for i := 2; i < 6; i++ {
		r, ok = Get(i)

		if !ok {
			t.Fatalf("expected record with id %v to exist. got %v, %v", i, r, ok)
		}
	}
}

func TestInsert(t *testing.T) {
	Init(3)

	for i := 0; i < 4; i++ { // add one more than the capacity to ensure we test the reassignment of the tail
		id := Insert(Record{Data: strconv.Itoa(i + 1)})

		if id != i+1 {
			t.Fatalf("expected assigned id to be sequential. got %v for loop idx of %v. expected %v", id, i, i+1)
		}
	}

	n, i, count := newest, 4, 0

	for n != nil {
		if strconv.Itoa(i) != n.r.Data {
			t.Fatalf("expected entry for %v when navigating from newest to oldest but got %v", i, n.r.Data)
		}
		i--
		count++
		n = n.older
	}

	if count != 3 {
		t.Fatalf("expected 3 nodes in the list but got %v", count)
	}

	n, i, count = oldest, 2, 3 // 1 should have been ejected so start at 2

	for n != nil {
		if strconv.Itoa(i) != n.r.Data {
			t.Fatalf("expected entry for %v navigating from oldest to newest but got %v", i, n.r.Data)
		}
		i++
		n = n.newer
	}

	if count != 3 {
		t.Fatalf("expected 3 nodes in the list but got %v", count)
	}
}

func BenchmarkInsert(b *testing.B) {
	Init(10000)

	b.ResetTimer()

	for i := 0; i < 1000000; i++ {
		Insert(Record{Data: "data"})
	}
}

func BenchmarkGet(b *testing.B) {
	Init(10000)

	for i := 0; i < 1000000; i++ {
		Insert(Record{Data: "data"})
	}

	b.ResetTimer()

	// contains IDs from 990001 to 1000000
	for i := 0; i < 1000; i++ {
		if _, ok := Get(1000000); !ok {
			b.Fatalf("missing upperbound data")
		}

		if _, ok := Get(990001); !ok {
			b.Fatalf("missing lowerbound data")
		}

		if _, ok := Get(995003); !ok {
			b.Fatalf("missing long seek from new to old data")
		}

		if _, ok := Get(994999); !ok {
			b.Fatalf("missing long seek from old to new data")
		}
	}
}
