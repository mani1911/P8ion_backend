package model

type Image struct {
	ID          uint   `gorm:"primarykey; not null"`
	UserID      uint   `gorm:"foreignKey; not null"`
	ImageBase64 string `gorm:"not null"`
	Content     string `gorm:"not null"`
}
