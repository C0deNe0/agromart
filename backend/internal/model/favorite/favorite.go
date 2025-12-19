package favorite

import (
	"github.com/C0deNe0/agromart/internal/model"
	"github.com/google/uuid"
)

type Favorite struct {
	model.Base

	UserID    uuid.UUID `json:"userId" db:"user_id"`
	ProductID uuid.UUID `json:"productId" db:"product_id"`
}
