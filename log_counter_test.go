package logger

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestCounter(t *testing.T) {
	counter := &Counter{}

	require.Equal(t, uint64(0), counter.Count())

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				fmt.Fprintln(counter, "write to count")
			}
		}()
	}
	wg.Wait()
	require.Equal(t, uint64(200), counter.Count())
	counter.Reset()
	require.Equal(t, uint64(0), counter.Count())
}
