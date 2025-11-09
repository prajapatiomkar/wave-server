package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"
)

type Hub struct {
	rooms          map[string]map[*Client]bool
	Broadcast      chan *IncomingMessage
	Register       chan *Client
	Unregister     chan *Client
	mu             sync.RWMutex
	messageHandler MessageHandler
}

type MessageHandler interface {
	HandleMessage(msg *IncomingMessage) (*OutgoingMessage, error)
}

func NewHub(messageHandler MessageHandler) *Hub {
	return &Hub{
		rooms:          make(map[string]map[*Client]bool),
		Broadcast:      make(chan *IncomingMessage),
		Register:       make(chan *Client),
		Unregister:     make(chan *Client),
		messageHandler: messageHandler,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			client.Hub = h

			h.mu.Lock()
			if h.rooms[client.RoomID] == nil {
				h.rooms[client.RoomID] = make(map[*Client]bool)
			}
			h.rooms[client.RoomID][client] = true
			h.mu.Unlock()

			log.Printf("✅ Client joined: %s (room: %s, total: %d)",
				client.Username, client.RoomID, len(h.rooms[client.RoomID]))

			h.broadcastToRoom(client.RoomID, &OutgoingMessage{
				Type:      "user_joined",
				Content:   client.Username + " joined the chat",
				UserID:    client.UserID,
				Username:  client.Username,
				RoomID:    client.RoomID,
				CreatedAt: time.Now(),
			}, nil)

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.rooms[client.RoomID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)

					if len(clients) == 0 {
						delete(h.rooms, client.RoomID)
					}
				}
			}
			h.mu.Unlock()

			log.Printf("❌ Client left: %s (room: %s)", client.Username, client.RoomID)

			h.broadcastToRoom(client.RoomID, &OutgoingMessage{
				Type:      "user_left",
				Content:   client.Username + " left the chat",
				UserID:    client.UserID,
				Username:  client.Username,
				RoomID:    client.RoomID,
				CreatedAt: time.Now(),
			}, nil)

		case message := <-h.Broadcast:
			outgoingMsg, err := h.messageHandler.HandleMessage(message)
			if err != nil {
				log.Printf("❌ Error handling message: %v", err)
				continue
			}

			h.broadcastToRoom(message.RoomID, outgoingMsg, nil)
		}
	}
}

func (h *Hub) broadcastToRoom(roomID string, message *OutgoingMessage, excludeClient *Client) {
	h.mu.RLock()
	clients := h.rooms[roomID]
	h.mu.RUnlock()

	if clients == nil {
		return
	}

	messageJSON, err := json.Marshal(message)
	if err != nil {
		log.Printf("Error marshaling message: %v", err)
		return
	}

	for client := range clients {
		if excludeClient != nil && client == excludeClient {
			continue
		}

		select {
		case client.Send <- messageJSON:
		default:
			close(client.Send)
			delete(clients, client)
		}
	}
}
