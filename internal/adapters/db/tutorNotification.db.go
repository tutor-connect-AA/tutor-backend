package db

import (
	"github.com/google/uuid"
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"gorm.io/gorm"
)

type TutorNotificationRepo struct {
	db *gorm.DB
}

func NewTutorNotificationRepo(db *gorm.DB) *TutorNotificationRepo {
	return &TutorNotificationRepo{
		db: db,
	}
}

type tutor_notification_table struct {
	gorm.Model
	Id          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OwnerId     uuid.UUID
	Tutor_table tutor_table `gorm:"foreignKey:OwnerId;references:Id"`
	Message     string
	Opened      bool
}

func (tNotf TutorNotificationRepo) CreateTutorNotification(newTutorNotification domain.Notification) (*domain.Notification, error) {
	ownerId, err := uuid.Parse(newTutorNotification.OwnerId)
	if err != nil {
		return nil, err
	}
	clientNotf := &tutor_notification_table{
		OwnerId: ownerId,
		Message: newTutorNotification.Message,
		Opened:  false,
	}
	if res := tNotf.db.Create(&clientNotf); res.Error != nil {
		return nil, res.Error
	}
	newTutorNotification.Id = clientNotf.Id.String()
	return &newTutorNotification, nil
}
func (tNotf TutorNotificationRepo) UpdateTutorNotificationStatus(id string) error {
	res := tNotf.db.Model(&tutor_notification_table{}).Where("id=?", id).UpdateColumn("opened", "true")
	return res.Error
}
func (tNotf TutorNotificationRepo) GetTutorNotificationById(id string) (*domain.Notification, error) {
	var tNtf tutor_notification_table
	if res := tNotf.db.First(&tNotf).Where("id=?", id); res.Error != nil {
		return nil, res.Error
	}
	return &domain.Notification{
		Id:        tNtf.Id.String(),
		Message:   tNtf.Message,
		OwnerId:   tNtf.OwnerId.String(),
		Opened:    tNtf.Opened,
		CreatedAt: tNtf.CreatedAt,
	}, nil

}
func (tNotf TutorNotificationRepo) GetTutorNotifications(ownerId string) ([]*domain.Notification, error) {
	var dTNtfs []*domain.Notification

	var tNtfs []tutor_notification_table

	res := tNotf.db.Order("created_at DESC").
		Where("owner_id = ?", ownerId).
		Find(&tNtfs)
	if res.Error != nil {
		return nil, res.Error
	}

	for _, tNtf := range tNtfs {
		dNtf := &domain.Notification{
			Id:        tNtf.Id.String(),
			Message:   tNtf.Message,
			OwnerId:   tNtf.OwnerId.String(),
			Opened:    tNtf.Opened,
			CreatedAt: tNtf.CreatedAt,
		}

		dTNtfs = append(dTNtfs, dNtf)
	}
	return dTNtfs, nil
}

func (tNotf TutorNotificationRepo) GetUnopenedTutorNotifications(ownerId string) ([]*domain.Notification, error) {
	var dTNtfs []*domain.Notification

	var tNtfs []tutor_notification_table

	if err := tNotf.db.Order("created_at DESC").Where("opened = ?", false).
		Where("owner_id = ?", ownerId).
		Find(&tNtfs).Error; err != nil {
		return nil, err
	}

	for _, tNtf := range tNtfs {
		dNtf := &domain.Notification{
			Id:        tNtf.Id.String(),
			Message:   tNtf.Message,
			OwnerId:   tNtf.OwnerId.String(),
			Opened:    tNtf.Opened,
			CreatedAt: tNtf.CreatedAt,
		}

		dTNtfs = append(dTNtfs, dNtf)
	}
	return dTNtfs, nil

}

func (tNotf TutorNotificationRepo) CountUnopenedTutorNotifications(ownerId string) (*int64, error) {

	var count int64
	if err := tNotf.db.Model(&tutor_notification_table{}).
		Where("owner_id = ?", ownerId).
		Where("opened = ?", false).
		Count(&count).Error; err != nil {

		return nil, err
	}
	return &count, nil

}

// CreateTutorNotification(newNotification domain.Notification) (*domain.Notification, error)
// UpdateTutorNotificationStatus(id string) error
// GetTutorNotificationById(id string) (*domain.Notification, error)
