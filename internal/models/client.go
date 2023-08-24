package models

import (
	"regexp"
)

var (
	ValidClientName = regexp.MustCompile(`^[a-z0-9_-]+$`)
)

type Client struct {
	Name string
}
