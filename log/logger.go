package log

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"path"
)

const ctxLoggerKey = "zapLogger"

var (
	Logger = &XLogger{
		zap.NewExample(),
	}

	ZapLevel = zap.NewAtomicLevelAt(zapcore.InfoLevel)
)

func InitLogger(logConf *LogConfig) *XLogger {

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "log",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	level, err := zapcore.ParseLevel(logConf.Level)
	if err != nil {
		fmt.Print("parse log level failed ", err)
		level = zapcore.InfoLevel
	}
	ZapLevel = zap.NewAtomicLevelAt(level)

	// 默认输出
	var writeSyncerList []zapcore.WriteSyncer
	if logConf.LogInConsole {
		writeSyncerList = append(writeSyncerList, os.Stdout)
	}
	writeSyncerList = append(writeSyncerList, zapcore.AddSync(createRotatingLogger(logConf, fmt.Sprintf("%s.%s.log", path.Join(logConf.Path, logConf.Name), level))))

	var coreList []zapcore.Core
	coreList = append(coreList, zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
		zapcore.NewMultiWriteSyncer(writeSyncerList...),
		ZapLevel, // 日志级别
	))

	// 每个级别的日志单独输出到一个文件中
	for l := level + 1; l <= zapcore.ErrorLevel; l++ {
		curLevel := l
		coreList = append(coreList, zapcore.NewCore(
			zapcore.NewJSONEncoder(encoderConfig), // 编码器配置
			zapcore.AddSync(createRotatingLogger(logConf, fmt.Sprintf("%s.%s.log", path.Join(logConf.Path, logConf.Name), curLevel.String()))),
			zap.LevelEnablerFunc(func(level zapcore.Level) bool {
				return level >= curLevel
			}), // 日志级别
		))

	}

	// 构造日志
	logger := zap.New(zapcore.NewTee(coreList...), zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	zap.ReplaceGlobals(logger)

	Logger = &XLogger{logger}
	Logger.Info("init log succ")
	return Logger
}

func createRotatingLogger(logConf *LogConfig, filename string) *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    logConf.MaxSize,
		MaxAge:     logConf.MaxAge,
		MaxBackups: logConf.MaxBackups,
		Compress:   logConf.Compress,
		LocalTime:  logConf.LocalTime,
	}
}

type XLogger struct {
	*zap.Logger
}

// WithValue Adds a field to the specified context
func (l *XLogger) WithValue(ctx context.Context, fields ...zapcore.Field) context.Context {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
		c.Request = c.Request.WithContext(context.WithValue(ctx, ctxLoggerKey, l.WithContext(ctx).With(fields...)))
		return c
	}
	return context.WithValue(ctx, ctxLoggerKey, l.WithContext(ctx).With(fields...))
}

// WithContext Returns a zap instance from the specified context
func (l *XLogger) WithContext(ctx context.Context) *XLogger {
	if c, ok := ctx.(*gin.Context); ok {
		ctx = c.Request.Context()
	}
	zl := ctx.Value(ctxLoggerKey)
	ctxLogger, ok := zl.(*zap.Logger)
	if ok {
		return &XLogger{ctxLogger}
	}
	return l
}
