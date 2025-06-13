package domain

import "github.com/google/uuid"

type Wallet struct {
	BaseModel
	Credits uint `gorm:"not null;default:10;check:credits >= 0"`
	UserID  uuid.UUID
}
