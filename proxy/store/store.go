// Store provides rolling storage of the last n entries using a fixed amount of memory relative to n
// It is implemented using a circular, doubly-linked list
package store

import (
	"comradequinn/hflow/log"
	"sync"
)

type (
	node struct {
		id    int
		r     Record
		older *node
		newer *node
	}
	Record struct {
		Data string
	}
)

var (
	cap    int
	length int
	newest *node
	oldest *node
	mx     sync.Mutex
)

func Init(capacity int) {
	mx.Lock()
	defer mx.Unlock()

	cap, length, newest, oldest = capacity, 0, nil, nil

	log.Printf(1, "store initialised with capacity of %v", capacity)
}

func Get(id int) (Record, bool) {
	mx.Lock()
	defer mx.Unlock()

	if id < oldest.id || id > newest.id {
		return Record{}, false
	}

	seekBackwards := newest.id-id < id-oldest.id // a minor optimisation to reduce search complexity from the naive O(n) to O(n/2)

	n := oldest

	if seekBackwards {
		n = newest
	}

	for n != nil {
		if n.id == id {
			return n.r, true
		}

		if seekBackwards {
			n = n.older
		} else {
			n = n.newer
		}
	}

	return Record{}, false
}

func Insert(r Record) int {
	mx.Lock()
	defer mx.Unlock()

	if cap <= 0 {
		panic("store.init not called")
	}

	id := 1

	if newest != nil {
		id = newest.id + 1
	}

	var n *node

	if length == cap {
		oldest = oldest.newer // set the second oldest as being the oldest
		n = oldest.older      // recycle the ejected node
		n.newer = nil
		oldest.older = nil // set the oldest as the end of the list
	} else {
		n = &node{}
		length++
	}

	n.id, n.r, n.older = id, r, newest

	if oldest == nil {
		oldest = n
	}

	if newest != nil {
		newest.newer = n
	}

	newest = n

	return id
}

func Len() int {
	mx.Lock()
	defer mx.Unlock()

	return length
}
