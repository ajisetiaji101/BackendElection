package model

import (
	"time"
)

// TimeOpenVote represents an entry in the time open vote table
type TimeOpenVote struct {
	ID           int       `db:"id"`
	OpenVoteTime string    `db:"date"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
