package model

import (
	"time"

	"github.com/google/uuid"
)

type BaseWithID struct {
	ID uuid.UUID `json:"id" db:"id"`
}

type BaseWithCreatedAt struct {
	CreatedAt time.Time `json:"created_at,omitempty" db:"created_at"`
}

type BaseWithUpdatedAt struct {
	UpdatedAt time.Time `json:"updated_at,omitempty" db:"updated_at"`
}

type Base struct {
	BaseWithID 
	BaseWithCreatedAt
	BaseWithUpdatedAt
}

type PaginatedResponse[T any] struct {
	Data       []T `json:"data"`
	Page       int `json:"page"`
	Limit      int `json:"limit"`
	Total      int `json:"total"`
	TotalPages int `json:"total_pages"`
}
