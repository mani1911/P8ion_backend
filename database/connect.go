package database

import (
	"fmt"
	"p8ion/config"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

func Connect() {
	config := config.GetConfig()
	dbPort := strconv.FormatUint(uint64(config.Db.Port), 10)

	connectionString := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local",
		config.Db.Username,
		config.Db.Password,
		config.Db.Host,
		dbPort,
		config.Db.Name,
	)

	var err error
	dblogger := gormlogger.Default.LogMode(gormlogger.Silent)

	// MySQL connection is established
	db, err = gorm.Open(mysql.Open(connectionString), &gorm.Config{Logger: dblogger})

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Database Connected Successfully")
	}
}
