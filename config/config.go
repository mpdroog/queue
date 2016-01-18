package config

import (
	"time"
	"os"
	"github.com/BurntSushi/toml"
)

type Server struct {
	Type string
	Hostname string
	Hashfeed string
	Connections int
	Queue int
	Timeout string
	SkipQueue int
}
type Queues struct {
	Steps []string
}
type Config struct {
	Servers map[string]Server
	Queues map[string]Queues
	ReplyTimeout string
}

var (
	C Config
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
	Timeout, e = time.ParseDuration(C.ReplyTimeout)
	if e != nil {
		return e
	}
	return nil
}