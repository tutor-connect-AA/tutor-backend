package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Adapter struct {
	db *gorm.DB
}

func NewAdapter(dsn string) (*Adapter, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})

	if err != nil {
		return nil, err
	}
	// err = db.AutoMigrate(&Order{}, OrderItem{})
	// if err != nil {
	// 	return nil, fmt.Errorf("db migration error: %v", err)
	// }

	return &Adapter{db: db}, nil
}
