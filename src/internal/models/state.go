package models

import (
	"time"
)

type State struct {
	Target Target
	RunAt  time.Time
}
