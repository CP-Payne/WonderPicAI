package domain

import (
	"time"

	"github.com/google/uuid"
)

type Status int

const (
	Failed Status = iota
	Pending
	PartiallyCompleted
	Completed
)

func (s Status) String() string {
	switch s {
	case Failed:
		return "Failed"
	case Pending:
		return "Pending"
	case Completed:
		return "Completed"
	case PartiallyCompleted:
		return "PartiallyCompleted"
	default:
		return "Unknown"
	}
}

type Prompt struct {
	BaseModel
	ExternalPromptID uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	Cost             int       `gorm:"not null"`
	ImageCount       int
	Width            int
	Height           int
	Images           []Image `gorm:"foreignKey:PromptID;references:ID"`
	Status           Status
	LastChecked      time.Time
}

type Image struct {
	BaseModel
	PromptID  uuid.UUID `gorm:"type:uuid;index;not null"`
	ImageData []byte    `gorm:"type:bytea"`
	Status    Status
}
