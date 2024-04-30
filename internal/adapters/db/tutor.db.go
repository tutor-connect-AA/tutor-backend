package db

import (
	"time"

	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
	"gorm.io/gorm"
)

type TutorRepo struct {
	db *gorm.DB
}

func NewTutorRepo(db *gorm.DB) *TutorRepo {
	return &TutorRepo{
		db: db,
	}
}

type tutor_table struct {
	gorm.Model
	Id                  string `gorm:"type:uuid;default:uuid_generate_v4()"`
	First_Name          string
	Fathers_Name        string
	Email               string `gorm:"unique"`
	Phone_Number        string
	Gender              domain.Gender
	Photo               string `gorm:"not null"`
	Rating              float32
	Bio                 string
	Username            string `gorm:"unique"`
	Password            string
	Role                Role
	CV                  string // Assuming CV is a file
	HourlyRate          float32
	Region              string
	City                string
	Education           domain.Education
	FieldOfStudy        string // Assuming limited options
	EducationCredential string // Assuming EducationCredential is a file
	CurrentlyEnrolled   domain.Education
	// ProofOfCurrentEnrollment string // Assuming ProofOfCurrentEnrollment is a file
	GraduationDate    time.Time
	PreferredSubjects string //pq.StringArray `gorm:"type:text[]"`   Limited to 2 choices
	// PreferredWorkLocation can be removed if not required for database storage
	// PreferredWorkLocation    string   // Optional
}

func (tr *TutorRepo) CreateTutorRepo(tutor *domain.Tutor) (*domain.Tutor, error) {
	hashedPass, err := utils.HashPass(tutor.Password)
	if err != nil {
		return nil, err
	}
	var newTutor = &tutor_table{
		Id:                  tutor.Id,
		First_Name:          tutor.FirstName,
		Fathers_Name:        tutor.FathersName,
		Email:               tutor.Email,
		Phone_Number:        tutor.PhoneNumber,
		Gender:              tutor.Gender,
		Photo:               tutor.Photo,
		Rating:              tutor.Rating,
		Bio:                 tutor.Bio,
		Username:            tutor.Username,
		Password:            hashedPass,
		Role:                Role(tutor.Role),
		CV:                  tutor.CV,
		HourlyRate:          tutor.HourlyRate,
		Region:              tutor.Region,
		City:                tutor.City,
		Education:           tutor.Education,
		FieldOfStudy:        tutor.FieldOfStudy,
		EducationCredential: tutor.EducationCredential,
		CurrentlyEnrolled:   tutor.CurrentlyEnrolled,
		// ProofOfCurrentEnrollment: tutor.ProofOfCurrentEnrollment,
		GraduationDate:    tutor.GraduationDate,
		PreferredSubjects: tutor.PreferredSubjects,
	}
	res := tr.db.Create(&newTutor)
	if res.Error != nil {
		return nil, res.Error
	}
	return tutor, nil
}

func (tr *TutorRepo) GetTutorByIdRepo(id string) (*domain.Tutor, error) {
	var tutor *tutor_table
	res := tr.db.Where("id = ?", id).First(&tutor)
	if res.Error != nil {
		return nil, res.Error
	}
	return &domain.Tutor{
		Id:                  tutor.Id,
		FirstName:           tutor.First_Name,
		FathersName:         tutor.Fathers_Name,
		Email:               tutor.Email,
		PhoneNumber:         tutor.Phone_Number,
		Gender:              tutor.Gender,
		Photo:               tutor.Photo,
		Rating:              tutor.Rating,
		Bio:                 tutor.Bio,
		Username:            tutor.Username,
		Password:            tutor.Password,
		Role:                domain.Role(tutor.Role),
		CV:                  tutor.CV,
		HourlyRate:          tutor.HourlyRate,
		Region:              tutor.Region,
		City:                tutor.City,
		Education:           tutor.Education,
		FieldOfStudy:        tutor.FieldOfStudy,
		EducationCredential: tutor.EducationCredential,
		CurrentlyEnrolled:   tutor.CurrentlyEnrolled,
		// ProofOfCurrentEnrollment: tutor.ProofOfCurrentEnrollment,
		GraduationDate:    tutor.GraduationDate,
		PreferredSubjects: tutor.PreferredSubjects,
	}, nil
}

func (tr *TutorRepo) GetTutorByEmail(email string) (*domain.Tutor, error) {
	var tutor *tutor_table
	res := tr.db.Where("username = ?", email).Find(&tutor)
	if res.Error != nil {
		return nil, res.Error
	}
	return &domain.Tutor{
		Id:                  tutor.Id,
		FirstName:           tutor.First_Name,
		FathersName:         tutor.Fathers_Name,
		Email:               tutor.Email,
		PhoneNumber:         tutor.Phone_Number,
		Gender:              tutor.Gender,
		Photo:               tutor.Photo,
		Rating:              tutor.Rating,
		Bio:                 tutor.Bio,
		Username:            tutor.Username,
		Password:            tutor.Password,
		Role:                domain.Role(tutor.Role),
		CV:                  tutor.CV,
		HourlyRate:          tutor.HourlyRate,
		Region:              tutor.Region,
		City:                tutor.City,
		Education:           tutor.Education,
		FieldOfStudy:        tutor.FieldOfStudy,
		EducationCredential: tutor.EducationCredential,
		CurrentlyEnrolled:   tutor.CurrentlyEnrolled,
		// ProofOfCurrentEnrollment: tutor.ProofOfCurrentEnrollment,
		GraduationDate: tutor.GraduationDate,
		// PreferredSubjects: tutor.PreferredSubjects,
	}, nil

}
