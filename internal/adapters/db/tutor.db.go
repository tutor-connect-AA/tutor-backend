package db

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"github.com/tutor-connect-AA/tutor-backend/internal/utils"
	"gorm.io/gorm"
)

// type User struct {
// 	db *gorm.DB
// }

// func NewTutorRepo(db *gorm.DB) *User {
// 	return &User{
// 		db: db,
// 	}
// }

type tutor_table struct {
	gorm.Model
	Id                  uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4()"`
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
	Role                domain.Role
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
	Applications  []job_application_table    `gorm:"foreignKey:applicant_id;references:Id"`
	Job_requests  []job_request_table        `gorm:"foreignKey:TutorId;references:Id"`
	Notifications []tutor_notification_table `gorm:"foreignKey:OwnerId;references:Id"`
}

func (ur *User) CreateTutorRepo(tutor *domain.Tutor) (*domain.Tutor, error) {
	hashedPass, err := utils.HashPass(tutor.Password)
	if err != nil {
		return nil, err
	}
	var newTutor = &tutor_table{
		// Id:                  tutor.Id,
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
		Role:                "TUTOR",
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

	newAuth := domain.Auth{
		Username: newTutor.Username,
		Password: newTutor.Password,
		Role:     newTutor.Role,
	}

	_, err = ur.CreateAuthRepo(newAuth)
	if err != nil {
		return nil, err
	}

	res := ur.db.Create(&newTutor)
	if res.Error != nil {
		return nil, res.Error
	}
	tutor.Id = newTutor.Id.String()
	return tutor, nil
}

func (ur *User) GetTutorByIdRepo(id string) (*domain.Tutor, error) {
	var tutor *tutor_table
	res := ur.db.Where("id = ?", id).First(&tutor)
	if res.Error != nil {
		return nil, res.Error
	}
	return &domain.Tutor{
		Id:                  tutor.Id.String(),
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

func (ur *User) GetTutorByUsername(username string) (*domain.Tutor, error) {
	var tutor *tutor_table
	res := ur.db.Where("username = ?", username).Find(&tutor)
	if res.Error != nil {
		if res.Error != nil {
			if res.Error == gorm.ErrRecordNotFound {
				return nil, domain.ErrNoRecord
			} else {
				return nil, res.Error
			}
		}
	}
	return &domain.Tutor{
		Id:                  tutor.Id.String(),
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
		Role:                tutor.Role,
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

func (ur *User) SearchTutorByNameRepo(name string) ([]*domain.Tutor, error) {

	var tutors []tutor_table
	res := ur.db.Raw("SELECT * FROM tutor_tables WHERE first_name ILIKE ? OR fathers_name ILIKE ?", "%"+name+"%", "%"+name+"%").Scan(&tutors)
	if res.Error != nil {
		return nil, res.Error
	}

	fmt.Println("Tutors from search at database is :", tutors)

	var tutorsDomain []*domain.Tutor

	for _, tutor := range tutors {
		tutDom := &domain.Tutor{
			Id:                  tutor.Id.String(),
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
			Role:                tutor.Role,
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
		}

		tutorsDomain = append(tutorsDomain, tutDom)
	}

	return tutorsDomain, nil

}

func (jr User) GetTutorsRepo(offset, limit int) ([]*domain.Tutor, error) {
	var tuts []tutor_table

	var tutList []*domain.Tutor

	if err := jr.db.Order("created_at").Offset(offset).Limit(limit).Find(&tuts).Error; err != nil {
		return nil, err
	}

	// res := jr.db.Find(&jbs)
	// if res.Error != nil {
	// 	return nil, res.Error
	// }

	for _, tutor := range tuts {
		oneTutor := &domain.Tutor{
			Id:                  tutor.Id.String(),
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
			Role:                tutor.Role,
			CV:                  tutor.CV,
			HourlyRate:          tutor.HourlyRate,
			Region:              tutor.Region,
			City:                tutor.City,
			Education:           tutor.Education,
			FieldOfStudy:        tutor.FieldOfStudy,
			EducationCredential: tutor.EducationCredential,
			CurrentlyEnrolled:   tutor.CurrentlyEnrolled,
			GraduationDate:      tutor.GraduationDate,
		}
		tutList = append(tutList, oneTutor)
	}
	return tutList, nil
}
