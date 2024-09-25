package client

import "context"

type Client interface {
	Subscribe(ctx context.Context, channel string, handler ClientHandler) error
}

type ClientHandler interface {
	Handle(data interface{})
}
