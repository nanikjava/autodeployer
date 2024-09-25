package client

import (
	"context"
	"log"

	"github.com/ably/ably-go/ably"
)

type RealAblyClient struct {
	client *ably.Realtime
	subscription func()

}

func NewRealAblyClient(apiKey string) (*RealAblyClient, error) {
	// Create a real Ably client
	client, err := ably.NewRealtime(ably.WithKey(apiKey))
	if err != nil {
		return nil, err
	}

	client.Connection.OnAll(func(change ably.ConnectionStateChange) {
		log.Printf("Connection event: %s state=%s reason=%s\n", change.Event, change.Current, change.Reason)
	})
	client.Connect()

	return &RealAblyClient{client: client}, nil
}


// Implement Subscribe method
func (c *RealAblyClient) Subscribe(ctx context.Context, channel string, handler func(msg *ably.Message)) error {
	sub, err := c.client.Channels.Get(channel).SubscribeAll(ctx, func(msg *ably.Message) {
		 handler(msg)
	})
	if err != nil {
		return err
	}


	c.subscription = sub
	return nil
}
