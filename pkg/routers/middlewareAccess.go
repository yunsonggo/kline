package routers

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/bitly/go-simplejson"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/yunsonggo/kline/pkg/xloger"
)

type AccessMiddleware struct {
	Logger  xloger.Logger
	SpanID  string
	TraceID string
}

func NewAccessMiddleware(logger xloger.Logger, spanID, traceID string) *AccessMiddleware {
	return &AccessMiddleware{
		Logger:  logger,
		SpanID:  spanID,
		TraceID: traceID,
	}
}

type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (am *AccessMiddleware) NewHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		am.ctxTrace(ctx)
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		method := ctx.Request.Method
		headers := am.flattenHeader(ctx.Request.Header)
		var rawData []byte
		if body, err := ctx.GetRawData(); err == nil {
			ctx.Request.Body = io.NopCloser(bytes.NewBuffer(body))
			rawData = body
		}

		rw := &responseWriter{
			ResponseWriter: ctx.Writer,
			body:           bytes.NewBufferString(""),
		}
		ctx.Writer = rw

		ctx.Next()

		respBody := rw.body.Bytes()
		end := time.Now()
		cost := fmt.Sprintf("%.4f", time.Since(start).Seconds())
		costValue := fmt.Sprintf("Start:%s,End:%s,Cost:%s", start.Format("2006-01-02 15:04:05"), end.Format("2006-01-02 15:04:05"), cost)

		requestBody, _ := am.passwordFilter(rawData)
		responseBody, _ := am.passwordFilter(respBody)

		am.Logger.Info(ctx, "[ACCESS]",
			xloger.Field{
				Key:   "PATH",
				Value: path,
			},
			xloger.Field{
				Key:   "METHOD",
				Value: method,
			},
			xloger.Field{
				Key:   "HEADERS",
				Value: headers,
			},
			xloger.Field{
				Key:   "QUERY",
				Value: query,
			},
			xloger.Field{
				Key:   "REQUEST",
				Value: string(requestBody),
			},
			xloger.Field{
				Key:   "RESPONSE",
				Value: string(responseBody),
			},
			xloger.Field{
				Key:   "COST",
				Value: costValue,
			},
		)
	}
}

func (am *AccessMiddleware) ctxTrace(ctx *gin.Context) (string, string) {
	var span, trace string
	if spanValue, ok := ctx.Get(am.SpanID); ok {
		span, _ = spanValue.(string)
	}
	if traceValue, ok := ctx.Get(am.TraceID); ok {
		trace, _ = traceValue.(string)
	}
	if span == "" {
		span = "Access"
		ctx.Set(am.SpanID, span)
	}
	if trace == "" {
		trace = uuid.New().String()
		ctx.Set(am.TraceID, trace)
	}
	return span, trace
}

// 将 Header 类型的 []string 值拼接成一个字符串
func (am *AccessMiddleware) flattenHeader(header http.Header) map[string]string {
	flattenedHeader := make(map[string]string)

	for key, values := range header {
		// 将 []string 拼接成一个字符串
		if key == "Authorization" {
			continue
		}
		flattenedHeader[key] = strings.Join(values, ", ")
	}

	return flattenedHeader
}

func (am *AccessMiddleware) passwordFilter(in []byte) ([]byte, error) {
	bodyJson, bodyErr := simplejson.NewJson(in)
	if bodyErr != nil {
		return in, bodyErr
	}
	password, err := bodyJson.Get("password").String()
	if err != nil {
		return in, err
	}
	if password == "" {
		return in, nil
	}
	bodyJson.Set("password", "")
	return bodyJson.Bytes()
}
