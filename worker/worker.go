package worker

import (
	"rqueue/config"
	"rqueue/queue"
	"fmt"
	"math/rand"
	"time"
)

func Init(s map[string]config.Server) error {
	for _, opts := range s {
		ch, e := queue.New(opts)
		if e != nil {
			return e
		}
		for i := 0; i < opts.Connections; i++ {
			go newConnection(ch)
		}
	}
	// TODO: validate hash not any edgecases with losing data?
	return nil
}

func newConnection(q chan queue.Job) {
	fmt.Printf("Span off worker\n")
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	for {
		// todo
		select {
			case job := <-q:
				fmt.Printf("Process %+v\n", job)
				// Done!
				rnd := r.Intn(3)
				if rnd == 0 {
					// Ok
					job.Update <- "DONE"
				} else if rnd == 1 {
					// err reply
					job.Update <- "ERR"
				} else if rnd == 2 {
					// timeout, do nothing and wait
				}
		}
	}
}