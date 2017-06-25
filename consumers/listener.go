package consumers

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/bsm/sarama-cluster"
)

//Handler that handler kafka messages received
type Handler func(ctx context.Context, msg *sarama.ConsumerMessage)

//Subscribe allow to attach an ahandler to a topic
func Subscribe(ctx context.Context, brokers []string, groupID, topicList string, offset int64, handler Handler) (*cluster.Consumer, error) {
	// Init config
	config := cluster.NewConfig()

	//config.Logger = logger //verbose mode
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = offset

	// Init consumer, consume errors & messages
	consumer, err := cluster.NewConsumer(brokers, groupID, strings.Split(topicList, ","), config)
	if err != nil {
		return nil, err
	}

	go func() {
		// Consume
		for {
			select {
			case msg, more := <-consumer.Messages():
				if more {
					handler(ctx, msg)
					consumer.MarkOffset(msg, "")
				}
			case ntf, more := <-consumer.Notifications():
				if more {
					fmt.Fprintf(os.Stdout, "Rebalanced: %+v\n", ntf)
				}
			case err, more := <-consumer.Errors():
				if more {
					fmt.Fprintf(os.Stderr, "Error: %s\n", err.Error())
				}
			}
		}
	}()

	return consumer, nil
}
