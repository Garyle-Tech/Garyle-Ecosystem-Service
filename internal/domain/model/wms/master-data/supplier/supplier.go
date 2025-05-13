package supplier

import (
	"database/sql"
	"time"
)

type Supplier struct {
	ID        int          `json:"id" db:"id"`
	Name      string       `json:"name" db:"name"`
	Address   string       `json:"address" db:"address"`
	Contact   string       `json:"contact" db:"contact"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"`
}

func (s *Supplier) IsDeleted() bool {
	return s.DeletedAt.Valid
}
