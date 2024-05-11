package db

import (
	"fmt"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"gorm.io/gorm"
)

type User struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) *User {
	return &User{
		db: db,
	}
}

type auth_table struct {
	gorm.Model
	Username string `gorm:"unique; not null"`
	Password string `gorm:"not null"`
	Role     domain.Role
}

func (aR *User) CreateAuthRepo(newAuth domain.Auth) error {
	auth := &auth_table{
		Username: newAuth.Username,
		Password: newAuth.Password,
		Role:     newAuth.Role,
	}

	res := aR.db.Create(&auth)
	return res.Error
}
func (aR *User) GetAuthByUsernameRepo(username string) (*domain.Auth, error) {
	var cred auth_table

	fmt.Println("username is : ", username)

	if res := aR.db.Where("username=?", username).First(&cred); res.Error != nil {
		return nil, res.Error
	}
	return &domain.Auth{
		Username: cred.Username,
		Password: cred.Password,
		Role:     cred.Role,
	}, nil
}
