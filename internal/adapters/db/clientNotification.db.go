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

func (cNotf ClientNotificationRepo) UpdateClientNotificationStatus(id string) error {
	res := cNotf.db.Model(&client_notification_table{}).Where("id=?", id).UpdateColumn("opened", "true")
	return res.Error
}
