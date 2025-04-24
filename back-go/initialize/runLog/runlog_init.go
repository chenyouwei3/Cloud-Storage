package runLog

import (
	conf "gin-web/initialize/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"time"
)

var ZapLog *zap.Logger

func InitRunLog() error {
	// 配置日志格式
	config := zapcore.EncoderConfig{
		TimeKey:      "time",
		LevelKey:     "level",
		MessageKey:   "msg",
		EncodeTime:   zapcore.ISO8601TimeEncoder,
		EncodeLevel:  zapcore.CapitalLevelEncoder,
		EncodeCaller: zapcore.ShortCallerEncoder,
	}

	// 日志文件路径（每天一个）
	today := time.Now().Format("2006-01-02")
	logFilePath := "./logs/" + today + ".log"

	// 确保 logs 目录存在
	if err := os.MkdirAll("./logs", os.ModePerm); err != nil {
		return err
	}

	// 文件存在就追加，不存在就创建
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	writeSyncer := zapcore.AddSync(file)

	// 创建 core
	var core zapcore.Core
	if conf.Conf.APP.Mode == "debug" {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(config),
			writeSyncer,
			zapcore.DebugLevel, // 开发环境建议是 Debug
		)
	} else {
		core = zapcore.NewCore(
			zapcore.NewJSONEncoder(config),
			writeSyncer,
			zapcore.InfoLevel,
		)
	}
	ZapLog = zap.New(core)
	// 不要在 Init 中 defer ZapLog.Sync()，应该在程序退出时调用
	// defer ZapLog.Sync() // 不推荐这样放这里

	return nil
}
