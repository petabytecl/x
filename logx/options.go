package logx

import (
	"io"
	"os"
	"time"

	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Options contains all possible settings
type Options struct {
	// Development configures the logger to use a Zap development config
	// (stacktraces on warnings, no sampling), otherwise a Zap production
	// config will be used (stacktraces on errors, sampling).
	Development bool
	// Encoder configures how Zap will encode the output.  Defaults to
	// console when Development is true and JSON otherwise
	Encoder zapcore.Encoder
	// DestWritter controls the destination of the log output.  Defaults to
	// os.Stderr.
	DestWritter io.Writer
	// Level configures the verbosity of the logging.  Defaults to Debug when
	// Development is true and Info otherwise
	Level zapcore.LevelEnabler
	// StacktraceLevel is the level at and above which stacktraces will
	// be recorded for all messages. Defaults to Warn when Development
	// is true and Error otherwise
	StacktraceLevel zapcore.LevelEnabler
	// ZapOpts allows passing arbitrary zap.Options to configure on the
	// underlying Zap logger.
	ZapOpts []zap.Option
}

// Opts allows to manipulate Options
type Opts func(*Options)

// BindFlags will parse the given flagset for zap option flags and set the log options accordingly
//  zap-devel: Development Mode defaults(encoder=consoleEncoder,logLevel=Debug,stackTraceLevel=Warn)
//			  Production Mode defaults(encoder=jsonEncoder,logLevel=Info,stackTraceLevel=Error)
//  log-encoder: Zap log encoding ('json' or 'console')
//  log-level:  Zap Level to configure the verbosity of logging. Can be one of 'debug', 'info', 'error',
//			       or any integer value > 0 which corresponds to custom debug levels of increasing verbosity")
//  stacktrace-level: Zap Level at and above which stacktraces are captured (one of 'warn' or 'error')
func (o *Options) BindFlags(fs *pflag.FlagSet) {

	// Set Development mode value
	fs.BoolVar(&o.Development, "debug", false, "Enable debug mode")

	// Set Encoder value
	var encVal encoderFlag
	encVal.setFunc = func(fromFlag zapcore.Encoder) {
		o.Encoder = fromFlag
	}
	fs.Var(&encVal, "log-encoder", "Log encoding ('json' or 'console')")

	// Set the Log Level
	var levelVal levelFlag
	levelVal.setFunc = func(fromFlag zapcore.LevelEnabler) {
		o.Level = fromFlag
	}
	fs.Var(&levelVal, "log-level",
		"Log level to configure the verbosity of logging. Can be one of 'debug', 'info', 'error', "+
			"or any integer value > 0 which corresponds to custom debug levels of increasing verbosity")

	// Set the StrackTrace Level
	var stackVal stackTraceFlag
	stackVal.setFunc = func(fromFlag zapcore.LevelEnabler) {
		o.StacktraceLevel = fromFlag
	}
	fs.Var(&stackVal, "stacktrace-level",
		"Log level at and above which stacktraces are captured (one of 'warn' or 'error')")
}

// UseFlagOptions configures the logger to use the Options set by parsing zap option flags from the CLI.
//  opts := zap.Options{}
//  opts.BindFlags(flag.CommandLine)
//  log := zap.New(zap.UseFlagOptions(&opts))
func UseFlagOptions(in *Options) Opts {
	return func(o *Options) {
		*o = *in
		o.addDefaults()
	}
}

// addDefaults adds defaults to the Options
func (o *Options) addDefaults() {
	if o.DestWritter == nil {
		o.DestWritter = os.Stderr
	}

	if o.Development {
		if o.Encoder == nil {
			encCfg := zap.NewDevelopmentEncoderConfig()
			o.Encoder = zapcore.NewConsoleEncoder(encCfg)
		}
		if o.Level == nil {
			lvl := zap.NewAtomicLevelAt(zap.DebugLevel)
			o.Level = &lvl
		}
		if o.StacktraceLevel == nil {
			lvl := zap.NewAtomicLevelAt(zap.WarnLevel)
			o.StacktraceLevel = &lvl
		}
		o.ZapOpts = append(o.ZapOpts, zap.Development())

	} else {
		if o.Encoder == nil {
			encCfg := zap.NewProductionEncoderConfig()
			o.Encoder = zapcore.NewJSONEncoder(encCfg)
		}
		if o.Level == nil {
			lvl := zap.NewAtomicLevelAt(zap.InfoLevel)
			o.Level = &lvl
		}
		if o.StacktraceLevel == nil {
			lvl := zap.NewAtomicLevelAt(zap.ErrorLevel)
			o.StacktraceLevel = &lvl
		}
		o.ZapOpts = append(o.ZapOpts,
			zap.WrapCore(func(core zapcore.Core) zapcore.Core {
				return zapcore.NewSamplerWithOptions(core, time.Second, 100, 100)
			}))
	}

	o.ZapOpts = append(o.ZapOpts, zap.AddStacktrace(o.StacktraceLevel))
}
