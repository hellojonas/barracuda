package logs

import (
	"bufio"
	"bytes"
	"sync"
	"testing"
)

func TestSholdLogNLines(t *testing.T) {
	buf := bytes.NewBuffer(make([]byte, 0))

	logger := New(buf)
	defer logger.Close()
	var wg sync.WaitGroup

	groups := 10
	groupLines := 1000
	for i := 0; i < groups; i++ {
		wg.Add(1)
		go func(id int) {
			for j := 0; j < groupLines; j++ {
				logger.Info("logger_id: %2d -> This is a info entry", id)
			}
			wg.Done()
		}(i)
	}

	wg.Wait()
	scanner := bufio.NewScanner(buf)
	lines := 0

	for scanner.Scan() {
		_ = scanner.Text()
		lines++
	}

	want := groups * groupLines

	if want != lines {
		t.Fatalf("want %d lines, got %d lines", want, lines)
	}
}
