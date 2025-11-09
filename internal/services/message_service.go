package services

import (
	"errors"
	"time"

	"github.com/prajapatiomkar/wave-server/internal/models"
	"github.com/prajapatiomkar/wave-server/internal/repositories"
	"github.com/prajapatiomkar/wave-server/internal/websocket"
)

type MessageService struct {
	messageRepo *repositories.MessageRepository
	userRepo    *repositories.UserRepository
}

func NewMessageService(messageRepo *repositories.MessageRepository, userRepo *repositories.UserRepository) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		userRepo:    userRepo,
	}
}

func (s *MessageService) HandleMessage(msg *websocket.IncomingMessage) (*websocket.OutgoingMessage, error) {
	if msg.Type == "typing" {
		return &websocket.OutgoingMessage{
			Type:      "typing",
			Content:   msg.Content,
			RoomID:    msg.RoomID,
			UserID:    msg.UserID,
			Username:  msg.Username,
			CreatedAt: time.Now(),
		}, nil
	}

	message := &models.Message{
		RoomID:  msg.RoomID,
		UserID:  msg.UserID,
		Content: msg.Content,
		Type:    "text",
	}

	if err := s.messageRepo.Create(message); err != nil {
		return nil, errors.New("failed to save message")
	}

	user, err := s.userRepo.FindByID(msg.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	return &websocket.OutgoingMessage{
		ID:        message.ID,
		Type:      "message",
		Content:   message.Content,
		RoomID:    message.RoomID,
		UserID:    message.UserID,
		Username:  user.Username,
		Avatar:    user.Avatar,
		CreatedAt: message.CreatedAt,
	}, nil
}

func (s *MessageService) GetMessageHistory(roomID string, limit, offset int) ([]models.MessageResponse, error) {
	messages, err := s.messageRepo.GetByRoom(roomID, limit, offset)
	if err != nil {
		return nil, err
	}

	var response []models.MessageResponse
	for _, msg := range messages {
		response = append(response, models.MessageResponse{
			ID:        msg.ID,
			RoomID:    msg.RoomID,
			UserID:    msg.UserID,
			Content:   msg.Content,
			Type:      msg.Type,
			CreatedAt: msg.CreatedAt,
			User: models.UserResponse{
				ID:       msg.User.ID,
				Username: msg.User.Username,
				Email:    msg.User.Email,
				FullName: msg.User.FullName,
				Avatar:   msg.User.Avatar,
			},
		})
	}

	return response, nil
}
