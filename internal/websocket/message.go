package websocket

import "time"

type IncomingMessage struct {
	Type     string `json:"type"`
	Content  string `json:"content"`
	RoomID   string `json:"room_id"`
	UserID   uint   `json:"user_id"`
	Username string `json:"username"`
}

type OutgoingMessage struct {
	ID        uint      `json:"id,omitempty"`
	Type      string    `json:"type"`
	Content   string    `json:"content"`
	RoomID    string    `json:"room_id"`
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	Avatar    string    `json:"avatar,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
