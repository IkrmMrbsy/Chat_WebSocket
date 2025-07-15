package main

import (
	"database/sql"
	"log"
	"net/http"
	"wschat/config"
	handler "wschat/internal/api/chat"
	"wschat/internal/service/chat"
	chatusecase "wschat/internal/usecase/chat"

	_ "github.com/lib/pq"
)

func main() {
	db, err := sql.Open(config.DBDriver, config.DBURL)
	if err != nil {
		log.Fatal(err)
	}

	chatRepo := chat.NewRepository(db)
	chatSvc := chat.NewService(chatRepo)
	chatUC := chatusecase.NewChatUsecase(chatSvc)
	chatHandler := handler.NewChatHandler(chatUC)

	http.HandleFunc("/ws", chatHandler.HandleWebSocket)
	http.HandleFunc("/api/chat/history", chatHandler.HandleWebSocket)

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
