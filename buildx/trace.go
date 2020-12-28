package buildx

import (
	"contrib.go.opencensus.io/exporter/jaeger"
	"github.com/opentracing/opentracing-go"
	"go.opentelemetry.io/otel/api/kv"
)

func (i Info) OpenTelemetryFields() []kv.KeyValue {
	return []kv.KeyValue{
		kv.String("version", i.Version),
		kv.String("revision", i.Revision),
		kv.String("branch", i.Branch),
		kv.String("build_user", i.BuildUser),
		kv.String("build_date", i.BuildDate),
		kv.String("go_version", i.GoVersion),
		kv.String("os", i.Os),
		kv.String("arch", i.Arch),
		kv.String("compiler", i.Compiler),
	}
}

func (i Info) OpenTracingFields() []opentracing.Tag {
	return []opentracing.Tag{
		{Key: "version", Value: i.Version},
		{Key: "revision", Value: i.Revision},
		{Key: "branch", Value: i.Branch},
		{Key: "build_user", Value: i.BuildUser},
		{Key: "build_date", Value: i.BuildDate},
		{Key: "go_version", Value: i.GoVersion},
		{Key: "os", Value: i.Os},
		{Key: "arch", Value: i.Arch},
		{Key: "compiler", Value: i.Compiler},
	}
}

func (i Info) OpenCensusFields() []jaeger.Tag {
	return []jaeger.Tag{
		jaeger.StringTag("version", i.Version),
		jaeger.StringTag("revision", i.Revision),
		jaeger.StringTag("branch", i.Branch),
		jaeger.StringTag("build_user", i.BuildUser),
		jaeger.StringTag("build_date", i.BuildDate),
		jaeger.StringTag("go_version", i.GoVersion),
		jaeger.StringTag("os", i.Os),
		jaeger.StringTag("arch", i.Arch),
		jaeger.StringTag("compiler", i.Compiler),
	}
}
