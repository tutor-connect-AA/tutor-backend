package db

import (
	"github.com/google/uuid"
	"github.com/tutor-connect-AA/tutor-backend/internal/application/core/domain"
	"gorm.io/gorm"
)

type ClientNotificationRepo struct {
	db *gorm.DB
}

func NewClientNotificationRepo(db *gorm.DB) *ClientNotificationRepo {
	return &ClientNotificationRepo{
		db: db,
	}
}

type client_notification_table struct {
	gorm.Model
	Id           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	OwnerId      uuid.UUID
	Client_Table client_table `gorm:"foreignKey:OwnerId;references:Id"`
	Message      string
	Opened       bool
}

func (cNotf ClientNotificationRepo) CreateClientNotification(newNotification domain.Notification) (*domain.Notification, error) {
	ownerId, err := uuid.Parse(newNotification.OwnerId)
	if err != nil {
		return nil, err
	}
	clientNotf := &client_notification_table{
		OwnerId: ownerId,
		Message: newNotification.Message,
		Opened:  false,
	}
	if res := cNotf.db.Create(&clientNotf); res.Error != nil {
		return nil, res.Error
	}
	return &newNotification, nil
}

func (cNotf *ClientNotificationRepo) UpdateClientNotificationStatus(id string) error {
	res := cNotf.db.Model(&client_notification_table{}).Where("id=?", id).UpdateColumn("opened", true)
	return res.Error
}

func (cNotf ClientNotificationRepo) GetClientNotificationById(id string) (*domain.Notification, error) {
	var cNtf client_notification_table
	if res := cNotf.db.First(&cNtf).Where("id=?", id); res.Error != nil {
		return nil, res.Error
	}
	return &domain.Notification{
		Id:        cNtf.Id.String(),
		Message:   cNtf.Message,
		OwnerId:   cNtf.OwnerId.String(),
		Opened:    cNtf.Opened,
		CreatedAt: cNtf.CreatedAt,
	}, nil
}

func (cNotf ClientNotificationRepo) GetClientNotifications(ownerId string) ([]*domain.Notification, error) {
	var dCNtfs []*domain.Notification

	var cNtfs []client_notification_table

	res := cNotf.db.Order("created_at DESC").Where("owner_id = ?", ownerId).Find(&cNtfs)
	if res.Error != nil {
		return nil, res.Error
	}

	for _, cNtf := range cNtfs {
		dNtf := &domain.Notification{
			Id:        cNtf.Id.String(),
			Message:   cNtf.Message,
			OwnerId:   cNtf.OwnerId.String(),
			Opened:    cNtf.Opened,
			CreatedAt: cNtf.CreatedAt,
		}

		dCNtfs = append(dCNtfs, dNtf)
	}
	return dCNtfs, nil
}

func (cNotf ClientNotificationRepo) GetUnopenedClientNotifications(ownerId string) ([]*domain.Notification, error) {
	var dTNtfs []*domain.Notification

	var tNtfs []client_notification_table

	if err := cNotf.db.Order("created_at DESC").Where("opened = ?", false).
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

func (cNotf ClientNotificationRepo) CountUnopenedClientNotifications(ownerId string) (*int64, error) {

	var count int64
	if err := cNotf.db.Model(&client_notification_table{}).
		Where("owner_id = ?", ownerId).
		Where("opened = ?", false).
		Count(&count).Error; err != nil {

		return nil, err
	}
	return &count, nil

}
