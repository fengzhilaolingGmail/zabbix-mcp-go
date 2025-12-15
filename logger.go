package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	// atomicLevel 支持运行时动态修改日志等级
	atomicLevel zap.AtomicLevel
)

// InitLogger 初始化日志记录器
func InitLogger(level string) error {
	// 创建logs目录
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 获取当前日期
	currentDate := time.Now().Format("2006-01-02")
	logFile := filepath.Join(logsDir, fmt.Sprintf("zabbix-mcp-%s.log", currentDate))

	// 配置日志编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:  "timestamp",
		LevelKey: "level",
		NameKey:  "logger",
		// 不记录 caller 字段
		CallerKey:      zapcore.OmitKey,
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "message",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		// 使用 FullCallerEncoder 输出日志打印的文件完整路径和行号
		EncodeCaller: zapcore.FullCallerEncoder,
	}

	// 创建JSON编码器
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)

	// 创建文件写入器
	fileWriter := zapcore.AddSync(&dateRotatingWriter{
		filename: logFile,
		file:     nil,
	})

	// 创建控制台写入器
	consoleWriter := zapcore.Lock(os.Stdout)

	// 解析并设置日志等级（兼容大小写），默认 info
	lvl, err := zapcore.ParseLevel(strings.ToLower(level))
	if err != nil {
		lvl = zap.InfoLevel
	}
	atomicLevel = zap.NewAtomicLevelAt(lvl)

	// 创建多写入器（同时写入文件和控制台），使用 atomicLevel 以支持运行时调整
	core := zapcore.NewTee(
		zapcore.NewCore(jsonEncoder, fileWriter, atomicLevel),
		zapcore.NewCore(jsonEncoder, consoleWriter, atomicLevel),
	)

	// 创建logger（不添加 caller 信息）
	logger = zap.New(core)
	sugar = logger.Sugar()

	return nil
}

// GetLogger 获取logger实例
func GetLogger() *zap.Logger {
	return logger
}

// GetSugar 获取sugar logger实例
func GetSugar() *zap.SugaredLogger {
	return sugar
}

// SetLogLevel 用于在运行时或初始化阶段设置日志等级，支持字符串例如 "debug", "info", "warn", "error"
func SetLogLevel(level string) error {
	lvl, err := zapcore.ParseLevel(strings.ToLower(level))
	if err != nil {
		return err
	}
	// 如果 atomicLevel 还未初始化（InitLogger 未调用），则先初始化一个 atomicLevel
	if atomicLevel == (zap.AtomicLevel{}) {
		atomicLevel = zap.NewAtomicLevelAt(lvl)
		return nil
	}
	atomicLevel.SetLevel(lvl)
	return nil
}

// GetLogLevel 返回当前日志等级的字符串表示
func GetLogLevel() string {
	if atomicLevel == (zap.AtomicLevel{}) {
		return ""
	}
	return atomicLevel.Level().String()
}

// Info 二次封装：记录 info 级别日志，传入 msg, err, data
func Info(msg string, err error, data interface{}) {
	if sugar == nil {
		return
	}
	if err != nil {
		// 将 error 作为字符串记录，避免在 JSON encoder 中出现不可序列化类型
		sugar.Infow(msg, "error", err.Error(), "data", data)
	} else {
		sugar.Infow(msg, "data", data)
	}
}

// Error 二次封装：记录 error 级别日志，传入 msg, err, data
func Error(msg string, err error, data interface{}) {
	if sugar == nil {
		return
	}
	if err != nil {
		sugar.Errorw(msg, "error", err.Error(), "data", data)
	} else {
		sugar.Errorw(msg, "data", data)
	}
}

// Warn 二次封装：记录 warn 级别日志
func Warn(msg string, err error, data interface{}) {
	if sugar == nil {
		return
	}
	if err != nil {
		sugar.Warnw(msg, "error", err.Error(), "data", data)
	} else {
		sugar.Warnw(msg, "data", data)
	}
}

// Debug 二次封装：记录 debug 级别日志
func Debug(msg string, err error, data interface{}) {
	if sugar == nil {
		return
	}
	if err != nil {
		sugar.Debugw(msg, "error", err.Error(), "data", data)
	} else {
		sugar.Debugw(msg, "data", data)
	}
}

// Sync 同步日志
func Sync() {
	if logger != nil {
		logger.Sync()
	}
}

// dateRotatingWriter 按日期切分的日志写入器
type dateRotatingWriter struct {
	filename string
	file     *os.File
	date     string
}

func (w *dateRotatingWriter) Write(p []byte) (n int, err error) {
	currentDate := time.Now().Format("2006-01-02")

	// 检查是否需要切换日志文件
	if w.date != currentDate || w.file == nil {
		if w.file != nil {
			w.file.Close()
		}

		// 创建新的日志文件
		newFilename := filepath.Join("logs", fmt.Sprintf("zabbix-mcp-%s.log", currentDate))
		w.file, err = os.OpenFile(newFilename, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return 0, err
		}
		w.date = currentDate
		w.filename = newFilename
	}

	return w.file.Write(p)
}

func (w *dateRotatingWriter) Sync() error {
	if w.file != nil {
		return w.file.Sync()
	}
	return nil
}
