package myService

import (
	"context"
	"fmt"
	"os"

	"github.com/Shopify/sarama"
	proto "github.com/golang/protobuf/proto"
)

type MyServiceTracker interface {
	Track(context.Context, *MyServiceMessage)
}

// MakeMyServiceTrackerEndpoint ...
func MakeMyServiceTrackerEndpoint(s MyServiceTracker) func(context.Context, *sarama.ConsumerMessage) {
	return func(ctx context.Context, msg *sarama.ConsumerMessage) {
		message := decodeMessage(msg.Value)

		s.Track(ctx, message)
	}

}

func decodeMessage(msg []byte) *MyServiceMessage {
	message := MyServiceMessage{}

	if err := proto.Unmarshal(msg, &message); err != nil {
		fmt.Fprintf(os.Stderr, "Program not able to deserialize myServiceMessage")
		return nil
	}

	return &message
}
