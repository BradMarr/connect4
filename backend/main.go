package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type Elo uint16
type UUID uuid.UUID

type Instance struct {
	hostUUID uuid.UUID
	elo      Elo
}

var instances []Instance

func spawnInstance(hostUUID uuid.UUID, hostELO Elo) {
	var newInstance Instance = Instance {
		hostUUID: hostUUID,
		elo: hostELO, 
	}
	instances = append(instances, newInstance)
}

func route(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		http.Redirect(w, r, "https://connect4.costellar.net", http.StatusSeeOther)
	case "/pair":
		w.Header().Set("Content-Type", "application/json")

		var uuidHeader string = r.Header.Get("UUID")
		if uuidHeader == "" {
			http.Error(w, "UUID header is missing", http.StatusBadRequest)
			return
		}
        var eloHeader string = r.Header.Get("ELO")
		if eloHeader == "" {
			http.Error(w, "ELO header is missing", http.StatusBadRequest)
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

		if (false) {
			
		} else {
            spawnInstance(hostUUID, hostELO)
        }

	default:
		http.NotFound(w, r)
	}
}

func main() {
	http.HandleFunc("/", route)

    fmt.Println(uuid.New().String())
	fmt.Println("Starting server on port 80...")
	err := http.ListenAndServe(":80", nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
