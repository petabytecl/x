package logx

import (
	"context"
	"strconv"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gocloud.dev/server/requestlog"
	zapadapter "logur.dev/adapter/zap"
	"logur.dev/logur"
	k8zap "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// LogFields alias for fields
type LogFields map[string]interface{}

// Logger logx interface
type Logger interface {
	// logur interface
	logur.Logger

	// logur context aware interface
	logur.LoggerContext

	// go-cloud interface
	requestlog.Logger

	// WithFields interface extensions
	WithFields(fields map[string]interface{}) Logger
}

type logger struct {
	logur logur.Logger
}

// Trace pass msg and fields to the inner logger
func (l *logger) Trace(msg string, fields ...map[string]interface{}) { l.logur.Trace(msg, fields...) }

// Debug pass msg and fields to the inner logger
func (l *logger) Debug(msg string, fields ...map[string]interface{}) { l.logur.Debug(msg, fields...) }

// Info pass msg and fields to the inner logger
func (l *logger) Info(msg string, fields ...map[string]interface{}) { l.logur.Info(msg, fields...) }

// Warn pass msg and fields to the inner logger
func (l *logger) Warn(msg string, fields ...map[string]interface{}) { l.logur.Warn(msg, fields...) }

// Error pass msg and fields to the inner logger
func (l *logger) Error(msg string, fields ...map[string]interface{}) { l.logur.Error(msg, fields...) }

// TraceContext pass msg and fields to the inner logger
func (l *logger) TraceContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Trace(msg, fields...)
}

// DebugContext pass msg and fields to the inner logger
func (l *logger) DebugContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Debug(msg, fields...)
}

// InfoContext pass msg and fields to the inner logger
func (l *logger) InfoContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Info(msg, fields...)
}

// WarnContext pass msg and fields to the inner logger
func (l *logger) WarnContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Warn(msg, fields...)
}

// ErrorContext pass msg and fields to the inner logger
func (l *logger) ErrorContext(_ context.Context, msg string, fields ...map[string]interface{}) {
	l.Error(msg, fields...)
}

// Log implements go-cloud request logger interface
func (l *logger) Log(entry *requestlog.Entry) {
	timeStamp := entry.ReceivedTime.Add(entry.Latency)
	l.WithFields(LogFields{
		"httpRequest": LogFields{
			"requestMethod": entry.RequestMethod,
			"requestUrl":    entry.RequestURL,
			"requestSize":   entry.RequestHeaderSize + entry.RequestBodySize,
			"status":        entry.Status,
			"responseSize":  entry.ResponseHeaderSize + entry.ResponseBodySize,
			"userAgent":     entry.UserAgent,
			"remoteIp":      entry.RemoteIP,
			"referer":       entry.Referer,
			"latency":       string(appendLatency(nil, entry.Latency)),
		},
		"timestamp": LogFields{
			"seconds": timeStamp.Unix(),
			"nanos":   timeStamp.Nanosecond(),
		},
		"logging.googleapis.com/trace":  entry.TraceID.String(),
		"logging.googleapis.com/spanId": entry.SpanID.String(),
	}).Info("httpRequest")
}

func appendLatency(b []byte, d time.Duration) []byte {
	// Parses format understood by google-fluentd (which is looser than the documented LogEntry format).
	// See the comment at https://github.com/GoogleCloudPlatform/fluent-plugin-google-cloud/blob/e2f60cdd1d97e79ffe4e91bdbf6bd84837f27fa5/lib/fluent/plugin/out_google_cloud.rb#L1539
	b = strconv.AppendFloat(b, d.Seconds(), 'f', 9, 64)
	b = append(b, 's')
	return b
}

// WithFields pass fields to the inner logger
func (l *logger) WithFields(fields map[string]interface{}) Logger {
	return &logger{
		logur: logur.WithFields(l.logur, fields),
	}
}

// NewLogger returns a Logger implementation
func NewLogger(opts ...Opts) Logger {
	o := &Options{}
	for _, opt := range opts {
		opt(o)
	}
	o.addDefaults()

	// this basically mimics New<type>Config, but with a custom sink
	sink := zapcore.AddSync(o.DestWritter)
	o.ZapOpts = append(o.ZapOpts, zap.AddCallerSkip(1), zap.ErrorOutput(sink))
	log := zap.New(zapcore.NewCore(&k8zap.KubeAwareEncoder{Encoder: o.Encoder, Verbose: o.Development}, sink, o.Level))

	return &logger{
		logur: zapadapter.New(
			log.WithOptions(o.ZapOpts...),
		),
	}
}
