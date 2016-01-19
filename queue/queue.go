package queue

import (
	"time"
	"crypto/md5"
	"rqueue/config"
	"fmt"
	"strings"
	"sync"
)

var (
	Q map[string][]Server
	wg sync.WaitGroup
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
	defer close(j.Update)

	done := false
	replyTimeout := time.After(config.C.Queues[j.Type.String()].ReplyTimeoutD)
	steps := config.C.Queues[j.Type.String()].Steps

	wg.Add(1)
	s.Added++
Processing:
	for i := 0; i < len(steps); i++ {
		var serverTimeout <- chan time.Time
		isQueued := false
		for _, group := range Q[steps[i]] {
			valid := false
			if group.Hashfeed != nil {
				if group.Hashfeed.Match(j.Hash) {
					valid = true
				}
			} else {
				valid = true
			}

			if valid {
				if availQueue(group) {
					fmt.Printf("[%s][%s] Try\n", msgid, group.Hostname)
					group.Queue <- j
					j.Server = i
					s.QueueAdded++

					serverTimeout = time.After(group.Timeout)

					// We used 'for' to 'randomly' pick a server
					// from the group
					isQueued = true
					break
				} else {
					fmt.Printf("[%s][%s] Queue too full\n", msgid, group.Hostname)
					s.Full++
				}
			}
		}

		if isQueued {
			// Now wait for a response
			select {
				case reply := <-j.Update:
					if reply == "DONE" {
						fmt.Printf("[%s] Finished\n", msgid)
						done = true
						s.Success++
						break Processing
					} else {
						fmt.Printf("[%s] Error: %s (trying on next)\n", msgid, reply)
						s.Error++
					}
				case <- serverTimeout:
					fmt.Printf("[%s] Server Timeout (trying on next)\n", msgid)
					s.Timeout++

				case <- replyTimeout:
					fmt.Printf("[%s] Reply timeout (dropping request)\n", msgid)
					break Processing
			}
		}

	}
	if !done {
		fmt.Printf("[%s] Client cmd lost\n", msgid)
		s.Lost++
	}
	s.Processed++

	wg.Done()
}

// Wait until all done
// useful in closing routine.
func Wait() {
	wg.Wait()
}