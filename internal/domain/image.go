package domain

import (
	"time"

	"github.com/google/uuid"
)

type GenerationStatus int

const (
	Failed GenerationStatus = iota
	Pending
	Completed
)

func (s GenerationStatus) String() string {
	switch s {
	case Failed:
		return "Failed"
	case Pending:
		return "Pending"
	case Completed:
		return "Completed"
	default:
		return "Unknown"
	}
}

type Prompt struct {
	BaseModel
	PromptID     uuid.UUID `gorm:"type:uuid;uniqueIndex;not null"`
	Cost         int       `gorm:"not null"`
	NumberImages int
	Width        int
	Height       int
	Images       []Image `gorm:"foreignKey:PromptID;references:PromptID"`
	Status       GenerationStatus
	LastChecked  time.Time
}

type Image struct {
	BaseModel
	PromptID  uuid.UUID `gorm:"type:uuid;index;not null"`
	ImageData []byte    `gorm:"type:bytea"`
}
