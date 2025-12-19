package category

import (
	"github.com/C0deNe0/agromart/internal/model"
	"github.com/google/uuid"
)

type Category struct {
	model.Base

	Name     string     `json:"name" db:"name"`
	Slug     string     `json:"slug" db:"slug"`
	ParentID *uuid.UUID `json:"parentId,omitempty" db:"parent_id"`
	IsActive bool       `json:"isActive" db:"is_active"`
}
