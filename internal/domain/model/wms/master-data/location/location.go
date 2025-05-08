package location

import (
	"database/sql"
	"time"
)

type Location struct {
	ID        int          `json:"id" db:"id"`
	Code      string       `json:"code" db:"code"`
	Zone      string       `json:"zone" db:"zone"`
	Type      string       `json:"type" db:"type"`         //rack/bin/area
	Capacity  float64      `json:"capacity" db:"capacity"` //1000
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"`
}

// IsDeleted checks if the location is deleted
func (l *Location) IsDeleted() bool {
	return l.DeletedAt.Valid
}
