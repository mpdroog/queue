package queue

type Stat struct {
	Added int
	QueueAdded int
	Processed int
	Success int
	Timeout int
	Error int
	Full int
	Lost int
}

var s Stat

func Stats() Stat {
	return s
}