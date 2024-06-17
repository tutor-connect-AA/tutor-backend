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
func (jr User) FilterTutorRepo(gender domain.Gender, rating, hourlyMin, hourlyMax int, city string, education domain.Education, fieldOfStudy string) ([]*domain.Tutor, error) {
	fmt.Println("Gender for filter at filter repo is : ", gender)
	fmt.Println("Education for filter at filter repo is : ", education)
	// query := `
	// 	SELECT * FROM tutor_tables
	// 	WHERE
	// 		($1::text IS NULL OR gender = $1)
	// 		AND ($2::int IS NULL OR rating >= $2)
	// 		AND ($3::int IS NULL OR hourly_rate >= $3)
	// 		AND ($4::int IS NULL OR hourly_rate <= $4)
	// 		AND ($5::text IS NULL OR city ILIKE $5)
	// 		AND ($6::text IS NULL OR education ILIKE $6)
	// 		AND ($7::text IS NULL OR field_of_study ILIKE $7)
	// `
	// 	query := `
	// 	SELECT * FROM tutor_tables
	// 	WHERE
	// 	  ($1 IS NULL OR gender = $1)
	// 	  AND ($2 IS NULL OR rating >= $2)
	// 	  AND ($3 IS NULL OR hourly_rate >= $3)
	// 	  AND ($4 IS NULL OR hourly_rate <= $4)
	// 	  AND ($5 IS NULL OR city ILIKE $5)
	// 	  AND ($6 IS NULL OR education ILIKE $6)
	// 	  AND ($7 IS NULL OR field_of_study ILIKE $7);
	// `
	query := `
	SELECT * FROM tutor_tables
	WHERE 
		(LOWER(gender) = LOWER($1) OR $1 = '')
		AND (rating >= $2 OR $2 = 0)
		AND (hourly_rate >= $3 OR $3 = 0)
		AND (hourly_rate <= $4 OR $4 = 0)
		AND (LOWER(city) LIKE LOWER($5) OR $5 = '')
		AND (LOWER(education) LIKE LOWER($6) OR $6 = '')
		AND (LOWER(field_of_study) LIKE LOWER($7) OR $7 = '')
`
	var filtered []tutor_table
	err := jr.db.Raw(query,
		gender,            // Gender filter
		rating,            // Rating filter
		hourlyMin,         // Minimum hourly rate filter
		hourlyMax,         // Maximum hourly rate filter
		"%"+city+"%",      // City filter
		"%"+education+"%", // Education filter
		"%"+fieldOfStudy+"%").Scan(&filtered).Error

	if err != nil {
		return nil, err
	}

	var tutors []*domain.Tutor

	for _, tutor := range filtered {

		tt := &domain.Tutor{
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

		tutors = append(tutors, tt)
	}
	return tutors, nil
}

func nullifyEmptyString(s string) interface{} {
	if s == "" {
		return ""
	}
	return s
}

func nullifyZeroValue(i int) interface{} {
	if i == 0 {
		return 0
	}
	return i
}
