package input

import (
	"bufio"
	"net-neko/pubsub"
	"os"
)

type stdin struct {
	pubsub *pubsub.PubSub[string]
}

func NewStdin() *stdin {
	s := stdin{pubsub.New[string]()}

	go func() {
		sc := bufio.NewScanner(os.Stdin)
		for sc.Scan() {
			txt := sc.Text()
			s.pubsub.Publish(txt)
		}
	}()

	return &s
}

func (s *stdin) SubScribe() *pubsub.Subscriber[string] {
	return s.pubsub.Subscribe()
}

func (s *stdin) CloseSub(sub *pubsub.Subscriber[string]) {
	s.pubsub.Close(sub)
}
