package main

import "time"

type sentence struct {
	Content string

	Arrived time.Time
	Sent    time.Time
}
