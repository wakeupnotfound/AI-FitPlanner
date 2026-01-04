package logger

import (
	"fmt"
	"github.com/ai-fitness-planner/backend/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path/filepath"
)

var Logger *zap.Logger

func InitLogger() error {
	logConfig := config.GlobalConfig.Log

	// 确保日志目录存在
	logDir := filepath.Dir(logConfig.Filename)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 配置日志切割
	lumberJackLogger := &lumberjack.Logger{
		Filename:   logConfig.Filename,
		MaxSize:    logConfig.MaxSize,
		MaxAge:     logConfig.MaxAge,
		MaxBackups: logConfig.MaxBackups,
		LocalTime:  true,
		Compress:   true,
	}

	// 获取日志级别
	level := getLogLevel(logConfig.Level)

	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 创建编码器
	var encoder zapcore.Encoder
	if config.GlobalConfig.App.Mode == "debug" {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 创建多输出端
	core := zapcore.NewTee(
		// 文件输出
		zapcore.NewCore(encoder, zapcore.AddSync(lumberJackLogger), level),
		// 控制台输出（仅在debug模式）
		zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			zapcore.AddSync(os.Stdout),
			level,
		),
	)

	// 创建Logger
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	// 替换全局logger
	zap.ReplaceGlobals(Logger)

	Logger.Info("日志系统初始化完成",
		zap.String("level", logConfig.Level),
		zap.String("filename", logConfig.Filename),
	)

	return nil
}

func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}

// 封装常用日志函数
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// 带错误信息的Error
func Errorf(msg string, err error, fields ...zap.Field) {
	fields = append(fields, zap.Error(err))
	Logger.Error(msg, fields...)
}
