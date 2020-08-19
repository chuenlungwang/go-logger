package logger

import (
	"io"
	"sync/atomic"
)

// Counter 是一个 io.Writer，用于对写入进行计数
type Counter struct {
	count uint64
}

func (c *Counter) incr() {
	atomic.AddUint64(&c.count, 1)
}

// Reset 对 Counter 进行重置
func (c *Counter) Reset() {
	atomic.StoreUint64(&c.count, 0)
}

// Count 返回当前的 count 值
func (c *Counter) Count() uint64 {
	return atomic.LoadUint64(&c.count)
}

func (c *Counter) Write(p []byte) (n int, err error) {
	c.incr()
	return len(p), nil
}

// LogCounter 返回一个 LogListener 用于对日志级别大于 it 的日志写入进行计数
func LogCounter(counter *Counter, it Threshold) LogListener {
	return func(t Threshold) io.Writer {
		if t < it {
			return nil
		}
		return counter
	}
}
