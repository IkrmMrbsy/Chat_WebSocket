package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"wschat/internal/usecase/chat"
	"wschat/internal/ws"

	"github.com/gorilla/websocket"
)

// Map semua koneksi WebSocket aktif
var clients = make(map[*websocket.Conn]int64)
var mu sync.Mutex

type ChatHandler struct {
	usecase chat.ChatUsecase
}

func NewChatHandler(usecase chat.ChatUsecase) *ChatHandler {
	return &ChatHandler{usecase}
}

func (h *ChatHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ws.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "WebSocket Upgrade failed", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil || userID == 0 {
		log.Println("Invalid user_id")
		return
	}

	// Simpan koneksi client
	mu.Lock()
	clients[conn] = userID
	mu.Unlock()

	defer func() {
		mu.Lock()
		delete(clients, conn)
		mu.Unlock()
	}()

	log.Printf("User %d connected", userID)

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		err = h.usecase.HanldeIncomingMessage(r.Context(), chat.CreateIn{
			UserId:  userID,
			Message: string(msg),
		})
		if err != nil {
			log.Println("failed to save message:", err)
			continue
		}

		broadcastMessage(userID, string(msg))
	}
}

// Fungsi untuk mengirim pesan ke semua koneksi yang aktif
func broadcastMessage(senderID int64, message string) {
	mu.Lock()
	defer mu.Unlock()

	for conn, uid := range clients {
		// bisa di-skip jika tidak ingin kirim ke pengirim
		_ = uid // kalau mau skip sender, tambahkan: if uid == senderID { continue }

		err := conn.WriteMessage(websocket.TextMessage, []byte(message))
		if err != nil {
			log.Println("broadcast error:", err)
			conn.Close()
			delete(clients, conn)
		}
	}
}

// Endpoint untuk ambil history chat
func (h *ChatHandler) GetChatHistory(w http.ResponseWriter, r *http.Request) {
	userIDStr := r.URL.Query().Get("user_id")
	userID, err := strconv.ParseInt(userIDStr, 10, 64)
	if err != nil {
		http.Error(w, "invalid user_id", http.StatusBadRequest)
		return
	}

	messages, err := h.usecase.GetHistory(r.Context(), userID)
	if err != nil {
		http.Error(w, "failed to get messages", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
