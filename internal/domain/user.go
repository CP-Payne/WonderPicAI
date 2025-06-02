package domain

type User struct {
	BaseModel
	Username string `gorm:"unique;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}
