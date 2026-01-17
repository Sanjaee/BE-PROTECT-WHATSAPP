package util

import (
	"log"
	"sync"
)

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	// Registered clients per user ID
	clients map[string]map[*Client]bool

	// Register requests from the clients
	Register chan *Client

	// Unregister requests from clients
	Unregister chan *Client

	// Mutex for thread-safe access to clients map
	mu sync.RWMutex
}

// NewHub creates a new Hub instance
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[string]map[*Client]bool),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			h.mu.Unlock()
			log.Printf("Client registered for user %s. Total clients for user: %d", client.UserID, len(h.clients[client.UserID]))

		case client := <-h.Unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.UserID]; ok {
				if _, exists := clients[client]; exists {
					delete(clients, client)
					close(client.Send)
					log.Printf("Client unregistered for user %s. Remaining clients: %d", client.UserID, len(clients))
				}
				if len(clients) == 0 {
					delete(h.clients, client.UserID)
				}
			}
			h.mu.Unlock()
		}
	}
}

// BroadcastToUser sends a message to all clients of a specific user
func (h *Hub) BroadcastToUser(userID string, message []byte) {
	h.mu.RLock()
	clients, ok := h.clients[userID]
	if !ok {
		h.mu.RUnlock()
		log.Printf("No clients found for user %s", userID)
		return
	}

	// Create a copy of clients to avoid holding lock while sending
	clientList := make([]*Client, 0, len(clients))
	for client := range clients {
		clientList = append(clientList, client)
	}
	h.mu.RUnlock()

	// Send message to all clients
	for _, client := range clientList {
		select {
		case client.Send <- message:
		default:
			// Client send channel is full, close connection
			close(client.Send)
			h.mu.Lock()
			if clients, ok := h.clients[client.UserID]; ok {
				delete(clients, client)
				if len(clients) == 0 {
					delete(h.clients, client.UserID)
				}
			}
			h.mu.Unlock()
		}
	}
}

// NotifyLogout notifies all clients of a user to logout
func (h *Hub) NotifyLogout(userID string) {
	message := []byte(`{"type":"logout","message":"Logged in from another device"}`)
	h.BroadcastToUser(userID, message)
	log.Printf("Logout notification sent to user %s", userID)
}
