package logger

import (
	"bytes"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"sync"
	"testing"
)

func TestThresholds(t *testing.T) {
	SetStdoutThreshold(LevelError)
	require.Equal(t, LevelError, StdoutThreshold())
	SetLogThreshold(LevelCritical)
	require.Equal(t, LevelCritical, LogThreshold())
	require.NotEqual(t, LevelCritical, StdoutThreshold())
	SetStdoutThreshold(LevelWarn)
	require.Equal(t, LevelWarn, StdoutThreshold())
}

func TestDefaultLogging(t *testing.T) {
	var outputBuf, logBuf bytes.Buffer

	defaultNotepad.outHandle = &outputBuf
	defaultNotepad.logHandle = &logBuf

	SetLogThreshold(LevelWarn)
	SetStdoutThreshold(LevelError)

	TRACE.Println("trace")
	DEBUG.Println("debug")
	INFO.Println("info")
	WARN.Println("warn")
	ERROR.Println("error")
	CRITICAL.Println("critical")
	FATAL.Println("fatal")

	outputStr := outputBuf.String()
	logStr := logBuf.String()

	require.Contains(t, outputStr, "fatal")
	require.Contains(t, outputStr, "critical")
	require.Contains(t, outputStr, "error")
	require.NotContains(t, outputStr, "warn")
	require.NotContains(t, outputStr, "info")
	require.NotContains(t, outputStr, "debug")
	require.NotContains(t, outputStr, "trace")

	require.Contains(t, logStr, "fatal")
	require.Contains(t, logStr, "critical")
	require.Contains(t, logStr, "error")
	require.Contains(t, logStr, "warn")
	require.NotContains(t, logStr, "info")
	require.NotContains(t, logStr, "debug")
	require.NotContains(t, logStr, "trace")
}

func TestLogCounter(t *testing.T) {
	assert := require.New(t)
	defaultNotepad.logHandle = ioutil.Discard
	defaultNotepad.outHandle = ioutil.Discard

	errorCounter := &Counter{}

	SetStdoutThreshold(LevelTrace)
	SetLogThreshold(LevelTrace)
	SetLogListeners(LogCounter(errorCounter, LevelError))

	FATAL.Println("fatal")
	CRITICAL.Println("critical")
	WARN.Println("a warning")
	WARN.Println("another warning")
	INFO.Println("info")
	DEBUG.Println("debug")

	assert.Equal(uint64(2), errorCounter.Count())

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				ERROR.Print("error")
				c := errorCounter.Count()
				assert.LessOrEqual(uint64(j), c)
			}
		}()
	}
	wg.Wait()
	assert.Equal(uint64(102), errorCounter.Count())
}
