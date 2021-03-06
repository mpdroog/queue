package config

import (
	"github.com/BurntSushi/toml"
	"os"
	"time"
)

type Server struct {
	Type        string
	Hostname    string
	Hashfeed    string
	Connections int
	Queue       int
	Timeout     string
	SkipQueue   int
}
type Queues struct {
	Steps         []string
	ReplyTimeout  string
	ReplyTimeoutD time.Duration
}
type Config struct {
	Servers map[string]Server
	Queues  map[string]Queues
}

var (
	C       Config
	Verbose bool
	Timeout time.Duration
)

func Open(path string) error {
	r, e := os.Open(path)
	if e != nil {
		return e
	}
	if _, e := toml.DecodeReader(r, &C); e != nil {
		return e
	}
	for id, queue := range C.Queues {
		timeout, e := time.ParseDuration(queue.ReplyTimeout)
		if e != nil {
			return e
		}
		queue.ReplyTimeoutD = timeout
		C.Queues[id] = queue
	}
	return nil
}
