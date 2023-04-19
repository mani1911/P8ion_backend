package model

import (
	"p8ion/database"
)

func MigrateDB() {
	db := database.GetDB()

	for _, model := range []interface{}{
		// Include models here to auto migrate
		User{},
		Image{},
	} {
		if err := db.AutoMigrate(&model); err != nil {
			panic(err)
		}
	}
}
