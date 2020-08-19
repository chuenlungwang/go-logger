package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
)

type Threshold int

func (t Threshold) String() string {
	return prefixes[t]
}

const (
	LevelTrace Threshold = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
	LevelFatal
)

var prefixes = map[Threshold]string{
	LevelTrace:    "TRACE",
	LevelDebug:    "DEBUG",
	LevelInfo:     "INFO",
	LevelWarn:     "WARN",
	LevelError:    "ERROR",
	LevelCritical: "CRITICAL",
	LevelFatal:    "FATAL",
}

// LogListener 可以给给定的日志级别日志输出进行事件侦听，这样可以用于测试或者类似情况。
// 需要注意的是任何级别的日志都会调用此函数，所以对于不感兴趣的日志级别请返回 nil 来忽略。
// 另外，即便配置中说明日志不会输出，也可以用侦听器来记录，比如 error 的个数。
type LogListener func(t Threshold) io.Writer

type Notepad struct {
	TRACE    *log.Logger
	DEBUG    *log.Logger
	INFO     *log.Logger
	WARN     *log.Logger
	ERROR    *log.Logger
	CRITICAL *log.Logger
	FATAL    *log.Logger

	LOG      *log.Logger
	FEEDBACK *Feedback

	loggers      [7]**log.Logger
	logHandle    io.Writer
	logThreshold Threshold
	outHandle    io.Writer
	outThreshold Threshold

	flags        int
	prefix       string
	loglisteners []LogListener
}

// NewNotepad 创建一个新的 Notepad
func NewNotepad(
	logHandle, outHandle io.Writer,
	logThreshold, outThreshold Threshold,
	prefix string, flags int,
	loglisteners ...LogListener,
) *Notepad {
	n := &Notepad{loglisteners: loglisteners}

	n.loggers = [...]**log.Logger{&n.TRACE, &n.DEBUG, &n.INFO, &n.WARN, &n.ERROR, &n.CRITICAL, &n.FATAL}
	n.logHandle = logHandle
	n.outHandle = outHandle
	n.logThreshold = logThreshold
	n.outThreshold = outThreshold
	n.flags = flags

	if len(prefix) != 0 {
		n.prefix = "[" + prefix + "] "
	}

	n.LOG = log.New(n.logHandle, "LOG:    ", n.flags)
	n.FEEDBACK = &Feedback{out: log.New(n.outHandle, "", 0), log: n.LOG}

	n.init()
	return n
}

// init 根据 notepad 的日志级别初始化各个日志器
func (n *Notepad) init() {
	outAndLogHandle := io.MultiWriter(n.outHandle, n.logHandle)
	for i, logger := range n.loggers {
		t := Threshold(i)
		prefix := n.prefix + t.String() + " "
		switch {
		case t >= n.logThreshold && t >= n.outThreshold:
			*logger = log.New(n.createLogWriters(outAndLogHandle, t), prefix, n.flags)
		case t >= n.logThreshold:
			*logger = log.New(n.createLogWriters(n.logHandle, t), prefix, n.flags)
		case t >= n.outThreshold:
			*logger = log.New(n.createLogWriters(n.outHandle, t), prefix, n.flags)
		default:
			*logger = log.New(n.createLogWriters(ioutil.Discard, t), prefix, n.flags)
		}
	}
}

func (n *Notepad) createLogWriters(handle io.Writer, t Threshold) io.Writer {
	if len(n.loglisteners) == 0 {
		return handle
	}

	writers := []io.Writer{handle}
	for _, l := range n.loglisteners {
		lwr := l(t)
		if lwr != nil {
			writers = append(writers, lwr)
		}
	}
	if len(writers) == 1 {
		return handle
	}
	return io.MultiWriter(writers...)
}

// SetStdoutOutput 重定向 notepad 的标准输出
func (n *Notepad) SetStdoutOutput(output io.Writer) {
	n.outHandle = output
	n.init()
}

// SetStdoutThreshold 改变 stdout 日志输出的级别
func (n *Notepad) SetStdoutThreshold(t Threshold) {
	n.outThreshold = t
	n.init()
}

// GetStdoutThreshold 获取 stdout 日志输出的级别
func (n *Notepad) GetStdoutThreshold() Threshold {
	return n.outThreshold
}

// SetLogOutput 设置 log 日志输出的文件
func (n *Notepad) SetLogOutput(output io.Writer) {
	n.logHandle = output
	n.init()
}

// SetLogThreshold 设置 log 日志输出的级别
func (n *Notepad) SetLogThreshold(t Threshold) {
	n.logThreshold = t
	n.init()
}

// GetLogThreshold 获取 log 日志输出的级别
func (n *Notepad) GetLogThreshold() Threshold {
	return n.logThreshold
}

// SetPrefix 改变 notepad 的输出前缀，前缀包含在中括号中
// 如果指定空的 prefix 将不输出 prefix 信息
func (n *Notepad) SetPrefix(prefix string) {
	if len(prefix) != 0 {
		n.prefix = "[" + prefix + "] "
	} else {
		n.prefix = ""
	}
	n.init()
}

// SetFlags 指定 logger 展示的 flag，查看 log 包获取更多信息
func (n *Notepad) SetFlags(flags int) {
	n.flags = flags
	n.init()
}

// SetLogListeners 设置新的日志监控器，将会替换掉原有的，所以需要自行保存
func (n *Notepad) SetLogListeners(l ...LogListener) {
	n.loglisteners = l
	n.init()
}

// 用于输出纯信息，同时将信息和额外的文件、行号信息记录到日志中
type Feedback struct {
	out *log.Logger
	log *log.Logger
}

func (fb *Feedback) Println(v ...interface{}) {
	fb.output(fmt.Sprintln(v...))
}

func (fb *Feedback) Print(v ...interface{}) {
	fb.output(fmt.Sprint(v...))
}

func (fb *Feedback) Printf(format string, v ...interface{}) {
	fb.output(fmt.Sprintf(format, v...))
}

func (fb *Feedback) output(s string) {
	if fb.log != nil {
		fb.log.Output(2, s)
	}
	if fb.out != nil {
		fb.out.Output(2, s)
	}
}
