package auth

import (
	"time"

	"github.com/C0deNe0/agromart/internal/model"
	"github.com/google/uuid"
)


//All fields are hidden from JSON.
type RefreshToken struct {
	model.Base
	UserID    uuid.UUID `json:"-" db:"user_id"`
	TokenHash string    `json:"-" db:"token_hash"`
	UserAgent string   `json:"-" db:"user_agent"`
	IPAddress string   `json:"-" db:"ip_address"`
	ExpiresAt time.Time `json:"-" db:"expires_at"`
	RevokedAt *time.Time `json:"-" db:"revoked_at"`
}



// Refresh Token

// Lifetime: 30 days
// Stored: HTTP-only cookie or secure storage
// Rotated on every use
// One token = one session/device