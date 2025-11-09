package repositories

import (
	"github.com/prajapatiomkar/wave-server/internal/models"
	"gorm.io/gorm"
)

type MessageRepository struct {
	db *gorm.DB
}

func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(message *models.Message) error {
	return r.db.Create(message).Error
}

func (r *MessageRepository) GetByRoom(roomID string, limit, offset int) ([]models.Message, error) {
	var messages []models.Message
	err := r.db.
		Preload("User").
		Where("room_id = ?", roomID).
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&messages).Error
	return messages, err
}

func (r *MessageRepository) GetByID(id uint) (*models.Message, error) {
	var message models.Message
	err := r.db.Preload("User").First(&message, id).Error
	return &message, err
}
