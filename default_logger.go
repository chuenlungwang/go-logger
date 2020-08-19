package logger

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

var (
	TRACE    *log.Logger
	DEBUG    *log.Logger
	INFO     *log.Logger
	WARN     *log.Logger
	ERROR    *log.Logger
	CRITICAL *log.Logger
	FATAL    *log.Logger

	LOG      *log.Logger
	FEEDBACK *Feedback

	defaultNotepad *Notepad
)

func init() {
	defaultNotepad = NewNotepad(ioutil.Discard, os.Stdout, LevelInfo, LevelTrace, "", log.LstdFlags)
	reloadDefaultNotepad()
}

func reloadDefaultNotepad() {
	TRACE = defaultNotepad.TRACE
	DEBUG = defaultNotepad.DEBUG
	INFO = defaultNotepad.INFO
	WARN = defaultNotepad.WARN
	ERROR = defaultNotepad.ERROR
	CRITICAL = defaultNotepad.CRITICAL
	FATAL = defaultNotepad.FATAL

	LOG = defaultNotepad.LOG
	FEEDBACK = defaultNotepad.FEEDBACK
}

// SetStdoutThreshold 设置默认 notepad 的 stdout 的日志级别，默认是 LevelTrace
func SetStdoutThreshold(t Threshold) {
	defaultNotepad.SetStdoutThreshold(t)
	reloadDefaultNotepad()
}

// SetStdoutOutput 对默认 notepad 的 stdout 进行重定向，默认是 stdout 本身
func SetStdoutOutput(output io.Writer) {
	defaultNotepad.SetStdoutOutput(output)
	reloadDefaultNotepad()
}

// SetLogThreshold 设置默认 notepad 的 log 日志输出级别，默认为 LevelInfo
func SetLogThreshold(t Threshold) {
	defaultNotepad.SetLogThreshold(t)
	reloadDefaultNotepad()
}

// SetLogOutput 修改默认 notepad 的 log 日志输出的位置，默认为 ioutil.Discard
func SetLogOutput(output io.Writer) {
	defaultNotepad.SetLogOutput(output)
	reloadDefaultNotepad()
}

// SetPrefix 给默认 notepad 设置日志前缀，前缀将包含在中括号中，如果设置为空字符串
// 将不输出前缀
func SetPrefix(prefix string) {
	defaultNotepad.SetPrefix(prefix)
	reloadDefaultNotepad()
}

// SetFlags 给默认 notepad 设置日志 flag，参考 log 包获取信息
func SetFlags(flags int) {
	defaultNotepad.SetFlags(flags)
	reloadDefaultNotepad()
}

// SetLogListeners 给默认 notepad 设置新的日志监听器
func SetLogListeners(l ...LogListener) {
	defaultNotepad.SetLogListeners(l...)
	reloadDefaultNotepad()
}

// LogThreshold 获取默认 notepad 的当前 log 的日志级别
func LogThreshold() Threshold {
	return defaultNotepad.logThreshold
}

// StdoutThreshold 获取默认 notepad 的当前 stdout 的日志级别
func StdoutThreshold() Threshold {
	return defaultNotepad.outThreshold
}
