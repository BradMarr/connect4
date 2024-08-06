package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Elo uint16
type UUID uuid.UUID

type Instance struct {
	hostUUID uuid.UUID
	hostELO      Elo
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var instances []Instance

func hostInstance(hostUUID uuid.UUID, hostELO Elo, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error opening websocket:", err)
		return
	}
	defer conn.Close()

	messageType, data, err := conn.ReadMessage()
	if err != nil {
		fmt.Println("Error receiving message:", err)
		return
	}
	if string(data) != "lobby made" {
		fmt.Println("Error verifying websocket")
		return
	}

	instances = append(instances, Instance {
		hostUUID: hostUUID, 
		hostELO: hostELO,
	})

	time.Sleep(time.Second) //await other player then continue
	if err := conn.WriteMessage(messageType, []byte("matched")); err != nil {
		fmt.Println("Error sending matched confimation:", err)
		return
	}	

	for {
		messageType, data, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error receiving message:", err)
			return
		}
		fmt.Println("Received message:", string(data))

		if err := conn.WriteMessage(messageType, data); err != nil {
			fmt.Println("Error sending message:", err)
			return
		}
	}
}

func route(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		http.Redirect(w, r, "https://connect4.costellar.net", http.StatusSeeOther)
	case "/connect":
		queryParams := r.URL.Query()
		uuidHeader := queryParams.Get("UUID")
		if uuidHeader == "" {
			http.Error(w, "UUID parameter is missing", http.StatusBadRequest)
			return
		}
		eloHeader := queryParams.Get("ELO")
		if eloHeader == "" {
			http.Error(w, "ELO parameter is missing", http.StatusBadRequest)
			return
		}

		hostUUID, err := uuid.Parse(uuidHeader)
		if err != nil {
			http.Error(w, "Invalid UUID format", http.StatusBadRequest)
			return
		}

		hostELOValue, err := strconv.ParseUint(eloHeader, 10, 16)
		if err != nil {
			http.Error(w, "Invalid ELO format", http.StatusBadRequest)
			return
		}
		hostELO := Elo(hostELOValue)

		hostInstance(hostUUID, hostELO, w, r)

	default:
		http.NotFound(w, r)
	}
}

func main() {
	http.HandleFunc("/", route)
	fmt.Println("Starting server on port 80...")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
