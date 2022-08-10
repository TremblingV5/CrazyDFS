package utils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var Lg *zap.Logger

// InitLogger 初始化Logger
func InitLogger() (err error) {
	writeSyncer := getLogWriter(
		TotalConf.Log.Filename,
		TotalConf.Log.MaxSize,
		TotalConf.Log.MaxBackups,
		TotalConf.Log.MaxAge,
	)
	encoder := getEncoder()
	var l = new(zapcore.Level)
	err = l.UnmarshalText([]byte(TotalConf.Log.Level))
	if err != nil {
		return
	}
	core := zapcore.NewCore(encoder, writeSyncer, l)

	Lg = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	zap.ReplaceGlobals(Lg) // 替换zap包中全局的logger实例，后续在其他包中只需使用zap.L()调用即可
	return
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeDuration = zapcore.SecondsDurationEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func getLogWriter(filename string, maxSize, maxBackup, maxAge int) zapcore.WriteSyncer {
	lumberJackLogger := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackup,
		MaxAge:     maxAge,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func WriteLog(logType string, msg string, values ...string) error {
	if len(values)%2 != 0 {
		panic("日志信息的键与值数量不匹配")
	}

	logInfo := make([]zapcore.Field, len(values)/2)
	for i := 0; i < len(values); i += 2 {
		logInfo[i/2] = zap.String(values[i], values[i+1])
	}

	switch logType {
	case "info":
		Lg.Info(msg, logInfo...)
	case "error":
		Lg.Error(msg, logInfo...)
	case "warn":
		Lg.Warn(msg, logInfo...)
	default:
		Lg.Info(msg, logInfo...)
	}
	return nil
}
