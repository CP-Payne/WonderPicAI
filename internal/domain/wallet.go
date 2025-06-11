package domain

import "github.com/google/uuid"

type Wallet struct {
	BaseModel
	Credits uint `gorm:"not null;default:100;check:credits >= 0"`
	UserID  uuid.UUID
}
