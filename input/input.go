package input

import (
	"net-neko/pubsub"
)

type Input interface {
	SubScribe() *pubsub.Subscriber[string]
	CloseSub(sub *pubsub.Subscriber[string])
}
