package product

import (
	"database/sql"
	"time"
)

type Product struct {
	ID          int          `json:"id" db:"id"`
	Sku         string       `json:"sku" db:"sku"` // unique
	Name        string       `json:"name" db:"name"`
	Description string       `json:"description" db:"description"`
	Unit        string       `json:"unit" db:"unit"`           //kg/pcs/box/lot
	Weight      float64      `json:"weight" db:"weight"`       //24.5
	Dimension   string       `json:"dimension" db:"dimension"` //100x50x20
	CreatedAt   time.Time    `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at" db:"updated_at"`
	DeletedAt   sql.NullTime `json:"deleted_at" db:"deleted_at"`
}

// is product deleted?
func (p *Product) IsDeleted() bool {
	return p.DeletedAt.Valid
}
