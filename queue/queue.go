package queue

import (
	"time"
	"crypto/md5"
	"rqueue/config"
	"fmt"
	"strings"
)

var (
	Q map[string][]Server
)

func init() {
	Q = make(map[string][]Server)
}

// Add server to queue (it's conns will pick jobs from the queue)
func New(opts config.Server) (chan Job, error) {
	timeout, e := time.ParseDuration(opts.Timeout)
	if e != nil {
		return nil, e
	}
	feed, e := parseFeed(opts.Hashfeed)
	if e != nil {
		return nil, e
	}

	q := make(chan Job, opts.Queue)
	Q[ opts.Type ] = append(Q[ opts.Type ], Server{
		Hostname: opts.Hostname,
		Hashfeed: feed,
		Timeout: timeout,
		SkipQueue: opts.SkipQueue,
		Queue: q,
	})
	return q, nil
}

func Test() error {
	// TODO: Check if config is ok to use?
	return nil
}

// Check if queue is available
func availQueue(s Server) bool {
	if cap(s.Queue) - len(s.Queue) <= s.SkipQueue {
		fmt.Printf(" Skip queue, too full\n")
		return false
	}
	return true
}

// Add job to queue for processing
func Add(cmd string, msgid string) {
	cmd = strings.ToUpper(cmd)
	hash := md5.New()
	hash.Write([]byte(msgid))

	j := Job{
		Added: time.Now(),
		Hash: hash.Sum(nil),
		Type: cmdType[cmd],
		Update: make(chan string, 1),
	}

	done := false
	replyTimeout := time.After(config.C.Queues[j.Type.String()].ReplyTimeoutD)
	steps := config.C.Queues[j.Type.String()].Steps

Processing:
	for i := 0; i < len(steps); i++ {
		var serverTimeout <- chan time.Time
		for _, group := range Q[steps[i]] {
			if group.Hashfeed != nil {
				if group.Hashfeed.Match(j.Hash) {
					if availQueue(group) {
						fmt.Printf("[%s][%s] Try\n", msgid, group.Hostname)
						group.Queue <- j
						j.Server = i

						serverTimeout = time.After(group.Timeout)

						// We used 'for' to 'randomly' pick a server
						// from the group
						break
					} else {
						fmt.Printf("[%s] Queue too full\n", msgid)
					}
				}
			} else {
				if availQueue(group) {
					fmt.Printf("[%s][%s] Try\n", msgid, group.Hostname)
					group.Queue <- j
					j.Server = i

					serverTimeout = time.After(group.Timeout)

					// We used 'for' to 'randomly' pick a server
					// from the group
					break
				} else {
					fmt.Printf("[%s] Queue too full\n", msgid)
				}
			}
		}

		// Now wait for a response
		select {
			case reply := <-j.Update:
				if reply == "DONE" {
					fmt.Printf("[%s] Finished\n", msgid)
					done = true
					break Processing
				} else {
					fmt.Printf("[%s] Error: %s (trying on next)\n", msgid, reply)
				}
			case <- serverTimeout:
				fmt.Printf("[%s] Server Timeout (trying on next)\n", msgid)
			case <- replyTimeout:
				fmt.Printf("[%s] Reply timeout (dropping request)\n", msgid)
				break Processing
		}

	}
	if !done {
		fmt.Printf("[%s] Client cmd lost\n", msgid)
	}

}
