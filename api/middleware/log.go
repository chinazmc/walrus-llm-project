package middleware

import (
	"bytes"
	"github.com/duke-git/lancet/v2/cryptor"
	"github.com/duke-git/lancet/v2/random"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"time"
	"walrus_llm_project/log"
)

func RequestLogMiddleware(logger *log.XLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// The configuration is initialized once per request
		uuid, err := random.UUIdV4()
		if err != nil {
			return
		}
		trace := cryptor.Md5String(uuid)
		logger.WithValue(ctx, zap.String("trace", trace))
		logger.WithValue(ctx, zap.String("method", ctx.Request.Method))
		logger.WithValue(ctx, zap.String("user-agent", ctx.Request.UserAgent()))
		logger.WithValue(ctx, zap.String("path", ctx.Request.URL.Path))
		logger.WithValue(ctx, zap.String("url", ctx.Request.URL.RawQuery))
		logger.WithValue(ctx, zap.String("ip", ctx.ClientIP()))

		if logger.Level() <= zap.DebugLevel {
			if ctx.Request.Body != nil {
				bodyBytes, _ := ctx.GetRawData()
				ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes)) // 关键点
				logger.WithValue(ctx, zap.String("request_params", string(bodyBytes)))
			}
		}

		logger.WithContext(ctx).Info("Request")
		ctx.Next()
	}
}
func ResponseLogMiddleware(logger *log.XLogger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var blw *bodyLogWriter
		if logger.Level() <= zap.DebugLevel {
			blw = &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
			ctx.Writer = blw
		}

		startTime := time.Now()
		ctx.Next()
		latencyTime := time.Since(startTime)

		fields := []zapcore.Field{
			zap.Duration("time", latencyTime),
			zap.Int("status", ctx.Writer.Status()),
		}

		if logger.Level() <= zap.DebugLevel {
			if blw != nil {
				fields = append(fields, zap.String("response_body", blw.body.String()))
			}
		}

		logger.WithContext(ctx).Info("Response", fields...)
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
