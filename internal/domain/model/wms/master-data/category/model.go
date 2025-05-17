package category

import (
	"database/sql"
	"time"
)

type Category struct {
	ID        int          `json:"id" db:"id"`
	Name      string       `json:"name" db:"name"`
	ParentID  *int         `json:"parent_id" db:"parent_id"`
	CreatedAt time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt sql.NullTime `json:"deleted_at" db:"deleted_at"`
}

// IsDeleted checks if the location is deleted
func (c *Category) IsDeleted() bool {
	return c.DeletedAt.Valid
}
