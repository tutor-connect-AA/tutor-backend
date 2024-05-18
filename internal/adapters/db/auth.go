package db

import (
	"fmt"

	"github.com/google/uuid"
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
	Id       uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
	Username string    `gorm:"unique; not null"`
	Password string    `gorm:"not null"`
	// ClientID uuid.UUID    `gorm:"type:uuid;"`
	// TutorID  uuid.UUID    `gorm:"type:uuid;"`
	// Client   client_table `gorm:"constraint:OnDelete:CASCADE;foreignKey:ClientID"`
	// Tutor    tutor_table  `gorm:"constraint:OnDelete:CASCADE;foreignKey:TutorID"`
	Role domain.Role
}

func (aR *User) CreateAuthRepo(newAuth domain.Auth) (string, error) {

	// cltId, err := uuid.Parse(newAuth.ClientID)
	// if err != nil {
	// 	return "", err
	// }
	// tutId, err := uuid.Parse(newAuth.TutorID)
	// if err != nil {
	// 	return "", err
	// }
	auth := &auth_table{
		Username: newAuth.Username,
		Password: newAuth.Password,
		// ClientID: cltId,
		// TutorID:  tutId,
		Role: newAuth.Role,
	}

	res := aR.db.Create(&auth)
	return auth.Id.String(), res.Error
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
		// ClientID: cred.ClientID.String(),
		// TutorID:  cred.TutorID.String(),
		Role: cred.Role,
	}, nil
}
