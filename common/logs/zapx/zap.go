package zapx

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go.uber.org/zap"
)

// 用于设置 zap 的调用堆栈跳过的偏移量。它在日志中指示调用堆栈时跳过一定数量的栈帧，以便日志能够正确地显示调用位置。
const callerSkipOffset = 3

type ZapWriter struct {
	logger *zap.Logger
}

func NewZapWriter(opts ...zap.Option) (logx.Writer, error) {
	opts = append(opts, zap.AddCallerSkip(callerSkipOffset))
	logger, err := zap.NewProduction(opts...)
	if err != nil {
		return nil, err
	}

	return &ZapWriter{
		logger: logger,
	}, nil
}

func (w *ZapWriter) Alert(v interface{}) {
	w.logger.Error(fmt.Sprint(v))
}

// Close 方法用于关闭 zap.Logger 并刷新所有缓冲区，以确保日志写入持久化。
func (w *ZapWriter) Close() error {
	return w.logger.Sync() //将所有挂起的日志写入文件（如果有的话），并且清理资源。
}

// Debug 方法用于记录日志的 Debug 级别。
func (w *ZapWriter) Debug(v interface{}, fields ...logx.LogField) {
	w.logger.Debug(fmt.Sprint(v), toZapFields(fields...)...)
}

// Error 方法用于记录日志的 Error 级别。
func (w *ZapWriter) Error(v interface{}, fields ...logx.LogField) {
	w.logger.Error(fmt.Sprint(v), toZapFields(fields...)...)
}

// Info 方法用于记录日志的 Info 级别，功能与 Debug 和 Error 类似，只是日志级别不同。
func (w *ZapWriter) Info(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toZapFields(fields...)...)
}

// Severe 方法用于记录日志的 Fatal 级别，通常表示程序遇到严重错误并将退出。
func (w *ZapWriter) Severe(v interface{}) {
	w.logger.Fatal(fmt.Sprint(v))
}

// Slow 方法用于记录日志的 Warn 级别，通常用于表示程序的某些操作执行得比较慢或可能导致潜在问题。
func (w *ZapWriter) Slow(v interface{}, fields ...logx.LogField) {
	w.logger.Warn(fmt.Sprint(v), toZapFields(fields...)...)
}

// Stack 方法用于记录日志的 Error 级别，同时附加堆栈信息，帮助开发者排查错误。
func (w *ZapWriter) Stack(v interface{}) {
	w.logger.Error(fmt.Sprint(v), zap.Stack("stack"))
}

// Stat 方法用于记录日志的 Info 级别，通常用于记录统计信息或一些不紧急的操作信息。
func (w *ZapWriter) Stat(v interface{}, fields ...logx.LogField) {
	w.logger.Info(fmt.Sprint(v), toZapFields(fields...)...)
}

func toZapFields(fields ...logx.LogField) []zap.Field {
	zapFields := make([]zap.Field, 0, len(fields))
	for _, f := range fields {
		zapFields = append(zapFields, zap.Any(f.Key, f.Value))
	}
	return zapFields
}
