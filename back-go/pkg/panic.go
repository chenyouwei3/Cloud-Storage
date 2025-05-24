package pkg

import (
	"fmt"
	"runtime"
)

// 调用 `recover()` 获取 panic 的值，并格式化返回错误信息，防止程序崩溃。
func Recover() (err error) {
	// `recover()` 用于捕获 `panic`，如果有 `panic` 发生，则 `r` 不为 `nil`
	if r := recover(); r != nil {
		// 创建一个 65535 字节大小的缓冲区用于存储堆栈信息
		buf := make([]byte, 65535)

		// `runtime.Stack(buf, false)` 获取当前 Goroutine 的堆栈信息
		// `false` 表示仅获取当前 Goroutine，而不是所有 Goroutine
		l := runtime.Stack(buf, false)

		// 组合 `panic` 信息和堆栈信息，并返回 `error`
		err = fmt.Errorf(fmt.Sprintf("%v: %s", r, buf[:l]))
	}
	return
}
