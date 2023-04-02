package model

type User struct {
	ID       uint `gorm:"primarykey"`
	Username string
}
