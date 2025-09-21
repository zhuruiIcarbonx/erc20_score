package middleware

import (
	"bytes"
	"io"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/zhuruiIcarbonx/erc20_score/erc20_score_backend/src/logger"
)

type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w BodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// RLog 请求响应日志打印处理
// RLog() 是一个中间件函数,用于记录HTTP请求和响应的详细日志
// 主要功能包括:
// 1. 记录请求的URL路径、查询参数
// 2. 记录请求体内容
// 3. 记录响应体内容
// 4. 记录请求处理时间
// 5. 记录请求/响应的各种元数据(状态码、方法、IP等)
// 6. 使用zap日志库将信息写入日志
func RLog() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取原始请求路径和查询参数(避免被其他中间件修改)
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery

		// 读取并保存请求体
		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		requestBody, _ := ioutil.ReadAll(tee)
		c.Request.Body = ioutil.NopCloser(&buf)
		bodyLogWriter := &BodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		// 记录开始时间
		start := time.Now()

		// 调用下一个处理器
		c.Next()

		// 获取响应体
		responseBody := bodyLogWriter.body.Bytes()
		if len(c.Errors) > 0 {
			// 如果有错误,记录错误信息
			for _, e := range c.Errors.Errors() {
				logger.Log.Error().Msg(e)
			}
		} else {
			// 计算处理时间
			latency := float64(time.Now().Sub(start).Nanoseconds() / 1000000.0)
			// 记录请求和响应的详细信息
			logger.Log.Info().
				Int("status", c.Writer.Status()).
				Str("method", c.Request.Method).
				Str("function", c.HandlerName()).
				Str("path", path).
				Str("query", query).
				Str("ip", c.ClientIP()).
				Str("user-agent", c.Request.UserAgent()).
				Str("token", c.Request.Header.Get("session_id")).
				Str("content-type", c.Request.Header.Get("Content-Type")).
				Float64("latency", latency).
				Str("request", string(requestBody)).
				Str("response", string(responseBody)).
				Msg("")
		}
	}
}
