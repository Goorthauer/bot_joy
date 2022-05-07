package model

import (
	"context"
	"fmt"
	"os"

	"contrib.go.opencensus.io/exporter/jaeger"
	"go.opencensus.io/trace"
)

type TracingService struct {
	exporter *jaeger.Exporter
}

func NewTrace() *TracingService {
	var (
		err            error
		tracingService *TracingService
	)
	if tracing := os.Getenv("TRACING"); tracing == "ON" {
		tracingService.exporter, err = jaeger.NewExporter(jaeger.Options{
			AgentEndpoint:     "localhost:6831",
			CollectorEndpoint: "http://localhost:14268/api/traces",
			Process: jaeger.Process{
				ServiceName: "bot_joy",
			},
		})
		trace.RegisterExporter(tracingService.exporter)
		trace.ApplyConfig(trace.Config{DefaultSampler: trace.AlwaysSample()})
	}
	if err != nil {
		fmt.Println("NewTrace err:", err)
	}
	return tracingService

}

type SpanLocal struct {
	*trace.Span
}

func (sl *SpanLocal) addInt(key string, attr int64) *SpanLocal {
	sl.AddAttributes(trace.Int64Attribute(key, attr))
	return sl
}
func (sl *SpanLocal) addString(key, attr string) *SpanLocal {
	sl.AddAttributes(trace.StringAttribute(key, attr))
	return sl
}

func NewSpan(ctx context.Context, msg string) *SpanLocal {
	_, span := trace.StartSpan(ctx, msg)
	span.SpanContext()
	return &SpanLocal{span}
}
