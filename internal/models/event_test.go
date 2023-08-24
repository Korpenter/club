package models

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/Korpenter/club/internal/utils"
)

func TestEventString(t *testing.T) {
	currentTime := time.Now()
	client := "johndoe"
	errMsg := errors.New("test error")

	tests := []struct {
		name     string
		event    Event
		expected string
	}{
		{"ClientSat",
			Event{Code: ClientSat, Timestamp: currentTime, ClientName: client, TableID: 1},
			fmt.Sprintf("%s %d %s %d", utils.Format(currentTime), ClientSat, client, 1),
		},
		{"EventError",
			Event{Code: EventError, Timestamp: currentTime, ErrorMsg: errMsg},
			fmt.Sprintf("%s %d %v", utils.Format(currentTime), EventError, errMsg),
		},
		{"ClientLeft",
			Event{Code: ClientLeft, Timestamp: currentTime, ClientName: client},
			fmt.Sprintf("%s %d %s", utils.Format(currentTime), ClientLeft, client),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.event.String()
			if result != tt.expected {
				t.Errorf("Expected: %s, got: %s", tt.expected, result)
			}
		})
	}
}
