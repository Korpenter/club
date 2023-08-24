package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/Korpenter/club/internal/utils"
)

var ErrInvalidEvent = errors.New("invalid event")

const (
	ClientArrived = iota + 1
	ClientSat
	ClientWaiting
	ClientLeft
	ClientForceLeft = iota + 7
	ClientSatFromQueue
	EventError
)

type Event struct {
	Code       int
	Timestamp  time.Time
	ClientName string
	TableID    int
	ErrorMsg   error
}

func (e *Event) String() string {
	switch e.Code {
	case ClientSatFromQueue, ClientSat:
		return fmt.Sprintf("%s %d %s %d", utils.Format(e.Timestamp), e.Code, e.ClientName, e.TableID)
	case EventError:
		return fmt.Sprintf("%s %d %s", utils.Format(e.Timestamp), e.Code, e.ErrorMsg.Error())
	default:
		return fmt.Sprintf("%s %d %s", utils.Format(e.Timestamp), e.Code, e.ClientName)
	}
}
