package init_db

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/logx"
	"strings"
)

type KafkaWriter struct { //KafkaWriter 是一个自定义的日志写入器（Writer）
	Pusher *kq.Pusher //kq.Pusher 是 go-zero 框架中的 Kafka 推送客户端
}

func NewKafkaWriter(puash *kq.Pusher) *KafkaWriter {
	return &KafkaWriter{
		Pusher: puash,
	}
}

// Write 方法是 KafkaWriter 的核心方法，它实现了 logx.Writer 接口中的 Write 方法。logx 在写日志时会调用这个方法。
func (w *KafkaWriter) Write(p []byte) (n int, err error) {
	//strings.TrimSpace(string(p)) 会去除日志消息的前后空格和换行符。
	if err := w.Pusher.Push(strings.TrimSpace(string(p))); err != nil {
		return 0, err
	}
	return len(p), nil
}

func LogXKafka() *logx.Writer {
	//创建了一个 Kafka 推送客户端
	//指定 Kafka 集群的地址，通常是一个或多个 Kafka 节点的地址和端口。在这里，指定的是本地 Kafka 实例 localhost:9092。
	//返回一个封装了 Kafka 推送逻辑的 KafkaWriter 实例。KafkaWriter 的 Write 方法会将日志消息推送到 Kafka 主题中
	pusher := kq.NewPusher([]string{"localhost:9092"}, "cloud_storage-log")
	defer pusher.Close()
	//logx.NewWriter(...) 将 KafkaWriter 封装为 logx.Writer，这样就可以使用 go-zero 提供的日志功能将日志写入 Kafka。
	writer := logx.NewWriter(NewKafkaWriter(pusher))
	return &writer
}
