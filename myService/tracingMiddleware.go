package myService

import (
	"context"

	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

// TracingLoginTrackerService ...
type TracingMyServiceTrackerService struct {
	Next MyServiceTracker
}

// NewTracingService ...
func NewTracingService(l MyServiceTracker) MyServiceTracker {
	return TracingMyServiceTrackerService{
		Next: l,
	}
}

// Track ...
func (s TracingMyServiceTrackerService) Track(ctx context.Context, msg *MyServiceMessage) {
	span, ctx := SpanTrace(ctx, msg, "Global")
	defer span.Finish()
	s.Next.Track(ctx, msg)
}

// SpanTrace ...
func SpanTrace(ctx context.Context, msg *MyServiceMessage, name string) (opentracing.Span, context.Context) {
	var span opentracing.Span
	if msg != nil {
		span, ctx = opentracing.StartSpanFromContext(ctx, "MyService::"+name)
		//defer span.Finish()
		span.SetTag("ID", msg.ID)
		span.SetTag("Name", msg.Name)
		span.LogFields(
			log.String("ID", msg.ID),
			log.String("Name", msg.Name))
	}

	return span, ctx
}
