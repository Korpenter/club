package service

import (
	"testing"
	"time"

	"github.com/Korpenter/club/internal/config"
	"github.com/Korpenter/club/internal/models"
	"github.com/Korpenter/club/internal/utils"
)

type MockStorage struct {
	exists         bool
	tableIsFree    bool
	errorToReturn  error
	dequeuedClient *models.Client
	AllTables      map[int]*models.Table
}

func (m *MockStorage) AddClient(name string) error {
	return m.errorToReturn
}

func (m *MockStorage) CheckFreeTables() bool {
	return m.tableIsFree
}

func (m *MockStorage) EnqueueClient(name string) error {
	return m.errorToReturn
}

func (m *MockStorage) DequeueClient() *models.Client {
	return m.dequeuedClient
}

func (m *MockStorage) RemoveClient(name string) {}

func (m *MockStorage) FreedTableByClient(name string, timeSat time.Time) int {
	return 1
}

func (m *MockStorage) ClientExists(name string) bool {
	return m.exists
}

func (m *MockStorage) SetClientTable(name string, tableID int, timeSat time.Time) error {
	return m.errorToReturn
}

func (m *MockStorage) KickAllClientsAndClearTables(kickTime time.Time) {}

func (m *MockStorage) ClearAllClients() []*models.Client {
	return nil
}

func (m *MockStorage) GetAllTables() map[int]*models.Table {
	return m.AllTables
}

func TestClientArrive(t *testing.T) {
	tests := []struct {
		name      string
		timestamp time.Time
		client    string
		mock      *MockStorage
		wantErr   error
	}{
		{
			name:      "Client arrives on time",
			timestamp: time.Now().Add(-1 * time.Hour),
			client:    "zzzz",
			mock:      &MockStorage{},
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{OpeningTime: time.Now().Add(-2 * time.Hour)}
			s := New(cfg, tt.mock)

			err := s.ClientArrive(tt.timestamp, tt.client)
			if err != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestClientSit(t *testing.T) {
	tests := []struct {
		name      string
		timestamp time.Time
		client    string
		tableID   int
		mock      *MockStorage
		wantErr   error
	}{
		{
			name:      "Known client sits at free table",
			timestamp: time.Now(),
			client:    "aohn",
			tableID:   1,
			mock:      &MockStorage{exists: true},
			wantErr:   nil,
		},
		{
			name:      "Unknown client tries to sit",
			timestamp: time.Now(),
			client:    "diman",
			tableID:   2,
			mock:      &MockStorage{},
			wantErr:   ErrClientUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			s := New(cfg, tt.mock)

			err := s.ClientSit(tt.timestamp, tt.client, tt.tableID)
			if err != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestClientWait(t *testing.T) {
	tests := []struct {
		name      string
		timestamp time.Time
		client    string
		mock      *MockStorage
		wantErr   error
	}{
		{
			name:      "Client waits due to no free tables",
			timestamp: time.Now(),
			client:    "zzzz",
			mock:      &MockStorage{tableIsFree: false},
			wantErr:   nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			s := New(cfg, tt.mock)

			err := s.ClientWait(tt.timestamp, tt.client)
			if err != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestClientLeave(t *testing.T) {
	tests := []struct {
		name      string
		timestamp time.Time
		client    string
		mock      *MockStorage
		wantErr   error
	}{
		{
			name:      "Known client leaves and a table becomes free",
			timestamp: time.Now(),
			client:    "bab",
			mock:      &MockStorage{exists: true, dequeuedClient: &models.Client{Name: "oppa"}},
			wantErr:   nil,
		},
		{
			name:      "Unknown client tries to leave",
			timestamp: time.Now(),
			client:    "bab",
			mock:      &MockStorage{},
			wantErr:   ErrClientUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &config.Config{}
			s := New(cfg, tt.mock)

			_, _, err := s.ClientLeave(tt.timestamp, tt.client)
			if err != tt.wantErr {
				t.Errorf("Expected error: %v, got: %v", tt.wantErr, err)
			}
		})
	}
}

func TestCalcProfits(t *testing.T) {
	cfg := &config.Config{
		HourlyRate: 10,
	}
	time1, _ := utils.Parse("5:45")
	time2, _ := utils.Parse("01:15")
	tests := []struct {
		name      string
		tables    map[int]*models.Table
		mock      *MockStorage
		wantTotal int
	}{
		{
			name: "Calculate profits for multiple tables",
			tables: map[int]*models.Table{
				1: {TotalTime: time1},
				2: {TotalTime: time2},
			},
			mock: &MockStorage{
				AllTables: map[int]*models.Table{
					1: {TotalTime: time1},
					2: {TotalTime: time2},
				},
			},
			wantTotal: (6 + 2) * 10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := New(cfg, tt.mock)

			profits := s.CalcProfits()
			total := 0
			for _, p := range profits {
				total += p.Sum
			}
			if total != tt.wantTotal {
				t.Errorf("Expected total profit: %v, got: %v", tt.wantTotal, total)
			}
		})
	}
}
