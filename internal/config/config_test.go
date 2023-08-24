package config

import (
	"bufio"
	"strings"
	"testing"
	"time"
)

func TestNewConfig(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedErr string
		expectedCfg *Config
	}{
		{
			name:        "valid input",
			input:       "5\n08:00 16:00\n10\n",
			expectedErr: "",
			expectedCfg: &Config{
				NumberOfTables: 5,
				OpeningTime:    time.Date(0, 1, 1, 8, 0, 0, 0, time.UTC),
				ClosingTime:    time.Date(0, 1, 1, 16, 0, 0, 0, time.UTC),
				HourlyRate:     10,
			},
		},
		{
			name:        "invalid number of tables",
			input:       "-3\n08:00 16:00\n10\n",
			expectedErr: "-3",
		},
		{
			name:        "invalid number of arguments in table line",
			input:       "2 3 5\n08:00 16:00\n10\n",
			expectedErr: "2 3 5",
		},
		{
			name:        "invalid opening time format",
			input:       "3\n0800 16:00\n10\n",
			expectedErr: "0800 16:00",
		},
		{
			name:        "invalid closing time format",
			input:       "3\n08:00 1600\n10\n",
			expectedErr: "08:00 1600",
		},
		{
			name:        "invalid hourly rate",
			input:       "3\n08:00 16:00\n-10\n",
			expectedErr: "-10",
		},
		{
			name:        "closing time before opening time",
			input:       "3\n16:00 08:00\n10\n",
			expectedErr: "16:00 08:00",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scanner := bufio.NewScanner(strings.NewReader(tt.input))
			cfg, err := NewConfig(scanner)

			if tt.expectedErr != "" {
				if err == nil {
					t.Fatalf("expected error %v, got nil", tt.expectedErr)
				}
				if !strings.Contains(err.Error(), tt.expectedErr) {
					t.Fatalf("expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}
			if err != nil {
				t.Fatalf("didn't expect error but got %v", err)
			}

			if cfg.NumberOfTables != tt.expectedCfg.NumberOfTables || !cfg.OpeningTime.Equal(tt.expectedCfg.OpeningTime) ||
				!cfg.ClosingTime.Equal(tt.expectedCfg.ClosingTime) || cfg.HourlyRate != tt.expectedCfg.HourlyRate {
				t.Fatalf("expected config %+v, got %+v", tt.expectedCfg, cfg)
			}
		})
	}
}
