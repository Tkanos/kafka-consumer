package myService

import (
	"context"
	"fmt"
	"os"
)

// LoggingMyServiceTrackerService ...
type LoggingMyServiceTrackerService struct {
	Next MyServiceTracker
}

// NewLoggingService ...
func NewLoggingService(s MyServiceTracker) MyServiceTracker {
	return LoggingMyServiceTrackerService{
		Next: s,
	}
}

// Track ...
func (s LoggingMyServiceTrackerService) Track(ctx context.Context, msg *MyServiceMessage) {
	if msg != nil {
		fmt.Fprintf(os.Stdout, "%s/%s\n", msg.ID, msg.Name)
	}

	s.Next.Track(ctx, msg)
}
