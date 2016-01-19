package queue

import (
	"time"
)

//go:generate stringer -type=QueueType
type QueueType int

func (q QueueType) String() string {
	if q == Storage {
		return "storage"
	}
	if q == Header {
		return "header"
	}
	if q == Post {
		return "post"
	}
	panic("Invalid QueueType")
}

const (
	_                 = iota
	Storage QueueType = iota
	Header  QueueType = iota
	Post    QueueType = iota
)

// convert cmd to type
var cmdType map[string]QueueType = map[string]QueueType{
	"ARTICLE": Storage,
}

type Job struct {
	Added time.Time
	Hash  []byte

	Type   QueueType
	Server int

	Update chan string
}
