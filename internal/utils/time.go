package utils

import (
	"time"
)

const timeLayout = "15:04"

func Parse(timestamp string) (time.Time, error) {
	parseed, err := time.Parse(timeLayout, timestamp)
	if err != nil {
		return time.Time{}, err
	}
	return parseed, nil
}

func Format(timestamp time.Time) string {
	format := timestamp.Format(timeLayout)
	return format
}
