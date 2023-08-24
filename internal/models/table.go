package models

import (
	"time"
)

type Table struct {
	Id        int
	Client    *Client
	ClientSat time.Time
	TotalTime time.Time
}
