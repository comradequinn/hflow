package syncio

import (
	"strconv"
	"strings"
	"sync"
	"testing"
)

func TestWriter(t *testing.T) {
	sb := strings.Builder{}
	start, wg, limit, w, template := make(chan bool), sync.WaitGroup{}, 1000, NewWriter(&sb), "data from iteration: "

	for i := 0; i < limit; i++ {
		wg.Add(1)

		go func(i int) {
			_ = <-start
			w.Write([]byte(template + strconv.Itoa(i) + "\n"))
			defer wg.Done()
		}(i)
	}

	close(start)
	wg.Wait()

	data := sb.String()

	for i := 0; i < limit; i++ {
		s := template + strconv.Itoa(i)

		if !strings.Contains(data, s+"\n") {
			t.Fatalf("did find expected string in data: [%v]", s)
		}
	}
}
