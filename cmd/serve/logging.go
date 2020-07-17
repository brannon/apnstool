package serve

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
)

type LoggingResponseWriter interface {
	http.ResponseWriter

	SetContextValue(name LogContextName, value string)
}

type LogContextName string

const (
	LogContextAppId       LogContextName = "appId"
	LogContextDeviceToken LogContextName = "deviceToken"
)

const CommonLogTimeFormat = "2/Jan/2006:15:04:05 -0700"

type loggingResponseWriter struct {
	contextValues map[LogContextName]string
	statusCode    int

	rw  http.ResponseWriter
	req *http.Request
}

func WithLogging(handler http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, req *http.Request) {
		lrw := &loggingResponseWriter{
			rw:  rw,
			req: req,
		}

		handler(lrw, req)
	}
}

func SetLoggingContextValue(rw http.ResponseWriter, name LogContextName, value string) {
	lrw, ok := rw.(*loggingResponseWriter)
	if ok {
		lrw.SetContextValue(name, value)
	}
}

func (rw *loggingResponseWriter) Header() http.Header {
	return rw.rw.Header()
}

func (rw *loggingResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.rw.WriteHeader(statusCode)

	appId := ""
	deviceToken := ""

	if rw.contextValues != nil {
		if contextValue, found := rw.contextValues[LogContextAppId]; found {
			appId = contextValue
		}
		if contextValue, found := rw.contextValues[LogContextDeviceToken]; found {
			deviceToken = contextValue
		}
	}

	rw.emitHttpLog(
		clientIP(rw.req),
		appId,
		deviceToken,
		time.Now(),
		rw.req.Method,
		rw.req.URL.Path,
		rw.req.Proto,
		rw.statusCode,
	)
}

func (rw *loggingResponseWriter) Write(data []byte) (int, error) {
	return rw.rw.Write(data)
}

func (rw *loggingResponseWriter) SetContextValue(name LogContextName, value string) {
	if rw.contextValues == nil {
		rw.contextValues = make(map[LogContextName]string)
	}
	rw.contextValues[name] = value
}

func (rw *loggingResponseWriter) emitHttpLog(
	clientIP string,
	appId string,
	deviceToken string,
	time time.Time,
	method string,
	path string,
	protocol string,
	status int,
) {
	fmt.Printf(
		"%s %s %s [%s] \"%s %s %s\" %d -\n",
		valueOrDash(clientIP),
		valueOrDash(hash(appId)),
		valueOrDash(hash(deviceToken)),
		time.Format(CommonLogTimeFormat),
		method,
		path,
		protocol,
		status,
	)
}

func valueOrDash(value string) string {
	if value != "" {
		return value
	}
	return "-"
}

func hash(value string) string {
	if value == "" {
		return ""
	}
	h := md5.Sum([]byte(value))
	return hex.EncodeToString(h[:])
}

func clientIP(req *http.Request) string {
	if req.RemoteAddr != "" {
		i := strings.LastIndex(req.RemoteAddr, ":")
		if i != -1 {
			return req.RemoteAddr[:i]
		}
	}
	return req.RemoteAddr
}
