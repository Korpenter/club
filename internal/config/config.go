package config

import (
	"bufio"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/Korpenter/club/internal/utils"
)

type Config struct {
	NumberOfTables int
	OpeningTime    time.Time
	ClosingTime    time.Time
	HourlyRate     int

	FileScanner *bufio.Scanner
}

func NewConfig(scanner *bufio.Scanner) (*Config, error) {
	cfg := &Config{
		FileScanner: scanner,
	}

	if cfg.FileScanner.Scan() {
		line := cfg.FileScanner.Text()

		num, err := strconv.Atoi(line)
		if err != nil || num < 1 {
			return nil, errors.New(line)
		}
		cfg.NumberOfTables = num
	}

	if cfg.FileScanner.Scan() {
		line := cfg.FileScanner.Text()
		times := strings.Split(line, " ")
		if len(times) != 2 {
			return nil, errors.New(line)
		}

		opening, err := utils.Parse(times[0])
		if err != nil {
			return nil, errors.New(line)
		}
		cfg.OpeningTime = opening

		closing, err := utils.Parse(times[1])
		if err != nil || !closing.After(opening) {
			return nil, errors.New(line)
		}
		cfg.ClosingTime = closing
	}

	if cfg.FileScanner.Scan() {
		line := cfg.FileScanner.Text()

		rate, err := strconv.Atoi(line)
		if err != nil || rate < 1 {
			return nil, errors.New(line)
		}
		cfg.HourlyRate = rate
	}
	return cfg, nil
}
