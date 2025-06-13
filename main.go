package main

import (
	"log"
	"os"
	"time"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	_ = godotenv.Load()                     // ignore if no .env
	key := os.Getenv("POLYGON_KEY")
	if key == "" {
		log.Fatal("POLYGON_KEY env var not set")
	}

	ws, _, err := websocket.DefaultDialer.Dial(
		"wss://socket.polygon.io/stocks", nil,
	)
	must(err)
	defer ws.Close()

	// auth
	must(ws.WriteJSON(map[string]string{
		"action": "auth", "params": key,
	}))
	// subscribe to SPY trades
	must(ws.WriteJSON(map[string]string{
		"action": "subscribe", "params": "T.SPY",
	}))

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			time.Sleep(time.Second) // back off, then exit or reconnect
			return
		}
		log.Println(string(msg))
	}
}
