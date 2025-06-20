package domain

type User struct {
	BaseModel
	Username string `gorm:"not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Wallet   Wallet
	Prompts  []Prompt `gorm:"foreignKey:UserID;references:ID"`
}
