package myService

import (
	"context"
)

// Service ...
type Service struct {
}

// NewService ...
func NewService() MyServiceTracker {
	s := Service{}
	return s
}

// Track ...
func (s Service) Track(ctx context.Context, msg *MyServiceMessage) {
	// PUT ALL YOUR LOGIC HERE
}
