package db

import (
	"fmt"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// type Adapter struct {
// 	db *gorm.DB
// }

func ConnectDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		fmt.Printf("Error connecting to DB %v", err)
		return nil, err
	}
	err = db.AutoMigrate(&client_table{}, &job_table{}, &tutor_table{})
	if err != nil {
		return nil, fmt.Errorf("db migration error: %v", err)
	}

	return db, nil
}
