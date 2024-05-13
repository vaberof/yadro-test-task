package eventhandler

import "time"

type Event struct {
	Time       time.Time
	Type       uint8
	ClientName string
	TableId    int
}
