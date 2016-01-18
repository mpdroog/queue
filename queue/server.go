package queue

import (
	"time"
)

type Server struct {
	Hostname string
	Hashfeed *ServerHashfeed

	Timeout time.Duration
	SkipQueue int

	Queue chan Job
}