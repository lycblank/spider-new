package notify

import "context"

type NotifyArg struct {
	Title string
	Contents []string
}

type Notify interface {
	Send(ctx context.Context, arg NotifyArg) error
}

