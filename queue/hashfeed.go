package queue

import (
	"fmt"
	"encoding/binary"
)

// https://github.com/mpdroog/hash
type ServerHashfeed struct {
	From uint32
	To uint32
	Base uint32
	Offset int
}

func (s *ServerHashfeed) Match(digest []byte) bool {
	begin := 16-s.Offset-4
	pos := binary.BigEndian.Uint32(digest[begin:])

	modulo := pos % s.Base +1
	if modulo >= s.From && modulo <= s.To {
		return true
	}
	return false
}

func parseFeed(feed string) (*ServerHashfeed, error) {
	if len(feed) == 0 {
		// set no feed
		return nil, nil
	}
	s := &ServerHashfeed{}

	n, e := fmt.Sscanf(feed, "%d-%d/%d:%d", &s.From, &s.To, &s.Base, &s.Offset)
	if e != nil {
		return nil, e
	}
	if n != 4 {
		return nil, fmt.Errorf("Failed parsing feed-str. n=%d", n)
	}
	return s, nil
}
