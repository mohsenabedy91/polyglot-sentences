package middlewares

import (
	"bytes"
	"github.com/gin-gonic/gin"
	"github.com/mohsenabedy91/polyglot-sentences/internal/core/config"
	"github.com/mohsenabedy91/polyglot-sentences/pkg/logger"
	"io"
	"strings"
	"time"
)

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r bodyLogWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

func (r bodyLogWriter) WriteString(s string) (int, error) {
	r.body.WriteString(s)
	return r.ResponseWriter.WriteString(s)
}

func DefaultStructuredLogger(log logger.Logger) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if strings.Contains(ctx.Request.URL.Path, "/swagger") {
			ctx.Next()
		} else {
			bodyLogWriter := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: ctx.Writer}
			start := time.Now()
			path := ctx.FullPath()
			raw := ctx.Request.URL.RawQuery
			bodyBytes, _ := io.ReadAll(ctx.Request.Body)
			err := ctx.Request.Body.Close()
			if err != nil {
				return
			}
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			ctx.Writer = bodyLogWriter

			ctx.Next()

			param := gin.LogFormatterParams{}
			param.TimeStamp = time.Now()
			param.Latency = time.Since(start)
			param.ClientIP = ctx.ClientIP()
			param.Method = ctx.Request.Method
			param.StatusCode = ctx.Writer.Status()
			param.ErrorMessage = ctx.Errors.ByType(gin.ErrorTypePrivate).String()
			param.BodySize = ctx.Writer.Size()
			if raw != "" {
				path = path + "?" + raw
			}
			param.Path = path

			headers := map[string]string{}
			headers[config.AppDeviceHeaderKey] = ctx.Request.Header.Get(config.AppDeviceHeaderKey)
			headers[config.AppVersionHeaderKey] = ctx.Request.Header.Get(config.AppVersionHeaderKey)

			keys := map[logger.ExtraKey]interface{}{}
			keys[logger.Path] = param.Path
			keys[logger.ClientIp] = param.ClientIP
			keys[logger.Method] = param.Method
			keys[logger.Latency] = param.Latency
			keys[logger.StatusCode] = param.StatusCode
			keys[logger.ErrorMessage] = param.ErrorMessage
			keys[logger.BodySize] = param.BodySize
			keys[logger.Headers] = headers
			keys[logger.RequestBody] = string(bodyBytes)
			keys[logger.ResponseBody] = bodyLogWriter.body.String()

			log.Info(logger.RequestResponse, logger.API, "Request", keys)
		}
	}
}
