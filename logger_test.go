package logger

import (
	"bytes"
	"fmt"
	"github.com/stretchr/testify/require"
	"io"
	"io/ioutil"
	"log"
	"testing"
)

func TestNotepad(t *testing.T) {
	var logHandle, outHandle bytes.Buffer
	errorCounter := &Counter{count: 0}
	notepad := NewNotepad(&logHandle, &outHandle, LevelCritical, LevelError, "TestNotePad", 0, LogCounter(errorCounter, LevelError))

	require.Equal(t, LevelError, notepad.outThreshold)
	require.Equal(t, LevelCritical, notepad.logThreshold)

	notepad.DEBUG.Print("Some debug")
	notepad.ERROR.Print("Some error")
	notepad.CRITICAL.Print("Some critical error")

	require.Contains(t, logHandle.String(), "[TestNotePad] CRITICAL Some critical error")
	require.NotContains(t, logHandle.String(), "Some error")
	require.NotContains(t, outHandle.String(), "Some debug")
	require.Contains(t, outHandle.String(), "[TestNotePad] ERROR Some error")

	require.Equal(t, uint64(2), errorCounter.Count())
}

func TestNotepadLogListener(t *testing.T) {
	require := require.New(t)

	var errorBuf, infoBuf bytes.Buffer

	errorCapture := func(t Threshold) io.Writer {
		if t != LevelError {
			return nil
		}
		return &errorBuf
	}

	infoCapture := func(t Threshold) io.Writer {
		if t != LevelInfo {
			return nil
		}
		return &infoBuf
	}

	notepad := NewNotepad(ioutil.Discard, ioutil.Discard, LevelCritical, LevelError, "TestNotePad", 0, errorCapture, infoCapture)

	notepad.DEBUG.Println("Some debug")
	notepad.INFO.Println("Some info")
	notepad.INFO.Println("Some more info")
	notepad.ERROR.Println("Some error")
	notepad.ERROR.Println("Some more error")
	notepad.CRITICAL.Println("Some critical")
	notepad.CRITICAL.Println("Some more critical")

	require.Equal(`[TestNotePad] ERROR Some error
[TestNotePad] ERROR Some more error
`, errorBuf.String())
	require.Equal(`[TestNotePad] INFO Some info
[TestNotePad] INFO Some more info
`, infoBuf.String())
}

func TestThreshold_String(t *testing.T) {
	require.Equal(t, LevelCritical.String(), "CRITICAL")
	require.Equal(t, LevelDebug.String(), "DEBUG")
}

func BenchmarkLogPrintOnlyToCouter(b *testing.B) {
	var logHandle, outHandle bytes.Buffer
	notepad := NewNotepad(&logHandle, &outHandle, LevelCritical, LevelCritical, "TestNotePad", 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		notepad.INFO.Print("test")
	}
}

func BenchmarkFprintString(b *testing.B) {
	var buf bytes.Buffer
	log := log.New(&buf, "", 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		log.Print("test")
	}
}

func BenchmarkLogString(b *testing.B) {
	var buf bytes.Buffer
	log := log.New(&buf, "", 0)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buf.Reset()
		s := fmt.Sprint("test")
		log.Output(2, s)
	}
}
