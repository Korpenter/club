package handler

import (
	"bufio"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/Korpenter/club/internal/config"
	"github.com/Korpenter/club/internal/models"
	"github.com/Korpenter/club/internal/service"
	"github.com/Korpenter/club/internal/utils"
)

type MockService struct {
	ArriveError error
	SitError    error
	WaitError   error
	LeaveError  error
	Dequeued    *models.Client
	FreedTable  int
	Kicked      []*models.Client
	Profits     []*models.Profit
}

func (m *MockService) ClientArrive(timestamp time.Time, name string) error {
	return m.ArriveError
}

func (m *MockService) ClientWait(timestamp time.Time, name string) error {
	return m.WaitError
}

func (m *MockService) ClientLeave(timestamp time.Time, name string) (*models.Client, int, error) {
	return m.Dequeued, m.FreedTable, m.LeaveError
}

func (m *MockService) ClientSit(timestamp time.Time, name string, tableID int) error {
	return m.SitError
}

func (m *MockService) KickClients(kickTime time.Time) []*models.Client {
	return m.Kicked
}
func (m *MockService) CalcProfits() []*models.Profit {
	return m.Profits
}

func TestFileHandler_ProcessEvents(t *testing.T) {
	cfg := &config.Config{}
	time, _ := utils.Parse("10:00")
	tests := []struct {
		name     string
		input    string
		expected []*models.Event
		err      string
		mock     *MockService
	}{
		{
			name:  "successful client arrive",
			input: "10:00 1 client\n",
			expected: []*models.Event{
				{
					Code:       models.ClientArrived,
					Timestamp:  time,
					ClientName: "client",
				},
			},
			mock: &MockService{ArriveError: nil},
		},
		{
			name:  "client arrive with error",
			input: "10:00 1 boba\n",
			expected: []*models.Event{
				{
					Code:       models.ClientArrived,
					Timestamp:  time,
					ClientName: "boba",
				},
				{
					Code:      models.EventError,
					Timestamp: time,
					ErrorMsg:  service.ErrYouShallNotPass,
				},
			},
			mock: &MockService{ArriveError: service.ErrYouShallNotPass},
		},
		{
			name:  "successful client left",
			input: "10:00 4 diman\n",
			expected: []*models.Event{
				{
					Code:       models.ClientLeft,
					Timestamp:  time,
					ClientName: "diman",
				},
				{
					Code:       models.ClientSat,
					Timestamp:  time,
					ClientName: "orel",
					TableID:    5,
				},
			},
			mock: &MockService{
				LeaveError: nil,
				Dequeued:   &models.Client{Name: "orel"},
				FreedTable: 5,
			},
		},
		{
			name:  "client left with ErrClientUnknown error",
			input: "10:00 4 diman\n",
			expected: []*models.Event{
				{
					Code:       models.ClientLeft,
					Timestamp:  time,
					ClientName: "diman",
				},
				{
					Code:      models.EventError,
					Timestamp: time,
					ErrorMsg:  service.ErrClientUnknown,
				},
			},
			mock: &MockService{LeaveError: service.ErrClientUnknown},
		},
		{
			name:  "successful client sit",
			input: "10:00 2 diman 5\n",
			expected: []*models.Event{
				{
					Code:       models.ClientSat,
					Timestamp:  time,
					ClientName: "diman",
					TableID:    5,
				},
			},
			mock: &MockService{SitError: nil},
		},
		{
			name:  "client sit with error",
			input: "10:00 2 diman 5\n",
			expected: []*models.Event{
				{
					Code:       models.ClientSat,
					Timestamp:  time,
					ClientName: "diman",
					TableID:    5,
				},
				{
					Code:      models.EventError,
					Timestamp: time,
					ErrorMsg:  errors.New("sit error"),
				},
			},
			mock: &MockService{SitError: errors.New("sit error")},
		},

		{
			name:  "successful client wait",
			input: "10:00 3 diman\n",
			expected: []*models.Event{
				{
					Code:       models.ClientWaiting,
					Timestamp:  time,
					ClientName: "diman",
				},
			},
			mock: &MockService{WaitError: nil},
		},
		{
			name:  "client wait with ErrICanWaitNoLonger error",
			input: "10:00 3 diman\n",
			expected: []*models.Event{
				{
					Code:       models.ClientWaiting,
					Timestamp:  time,
					ClientName: "diman",
				},
				{
					Code:      models.EventError,
					Timestamp: time,
					ErrorMsg:  service.ErrICanWaitNoLonger,
				},
			},
			mock: &MockService{WaitError: service.ErrICanWaitNoLonger},
		},
		{
			name:  "client wait with ErrQueueFull error",
			input: "10:00 3 diman\n",
			expected: []*models.Event{
				{
					Code:       models.ClientWaiting,
					Timestamp:  time,
					ClientName: "diman",
				},
				{
					Code:       models.ClientForceLeft,
					Timestamp:  time,
					ClientName: "diman",
				},
			},
			mock: &MockService{WaitError: service.ErrQueueFull},
		},
		{
			name:  "malformed event with missing client name",
			input: "10:00 2\n",
			err:   "10:00 2",
			mock:  &MockService{},
		},
		{
			name:  "malformed event with invalid time format",
			input: "10:70 2 diman\n",
			err:   "10:70 2 diman",
			mock:  &MockService{},
		},
		{
			name:  "malformed event with invalid event code",
			input: "10:00 x diman\n",
			err:   "10:00 x diman",
			mock:  &MockService{},
		},
		{
			name:  "malformed event with missing tableID for ClientSat event",
			input: "10:00 2 diman\n",
			err:   "10:00 2 diman",
			mock:  &MockService{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := strings.NewReader(tt.input)
			scanner := bufio.NewScanner(r)
			handler := NewFileHandler(scanner, tt.mock, cfg)
			err := handler.ProcessEvents()

			if tt.err != "" {
				if err == nil || err.Error() != tt.err {
					t.Errorf("Expected error: %s, got: %v", tt.err, err)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error, but got: %v", err)
				} else if !reflect.DeepEqual(tt.expected, handler.ee) {
					t.Errorf("Expected events: %v, got: %v", tt.expected, handler.ee)
				}
			}
		})
	}
}
