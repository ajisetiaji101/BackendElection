package model

import "time"

// IPWhitelist represents an entry in the IP whitelist table
type IPWhitelist struct {
	ID          int       `db:"id"`
	IPAddress   string    `db:"ip_address"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}
