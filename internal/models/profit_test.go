package models

import (
	"testing"
	"time"

	"github.com/Korpenter/club/internal/utils"
)

func TestProfitString(t *testing.T) {
	table := &Table{
		Id:        1,
		Client:    &Client{Name: "pippa"},
		ClientSat: time.Now(),
		TotalTime: time.Now().Add(1 * time.Hour),
	}

	tests := []struct {
		name     string
		profit   Profit
		expected string
	}{
		{
			"BasicProfitTest",
			Profit{Table: table, Sum: 100},
			"1 100 " + utils.Format(table.TotalTime),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.profit.String()
			if result != tt.expected {
				t.Errorf("Expected: %s, got: %s", tt.expected, result)
			}
		})
	}
}
