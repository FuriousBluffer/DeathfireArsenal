package api_handlers

import (
	"DeathfireArsenal/internal/errormanagement"
	"DeathfireArsenal/internal/logic"
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"net/http"
	"strings"
)

type APIHandlers struct {
	Logic *logic.BusinessLogic
}

func (a *APIHandlers) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	panic("implement me")
}

func (a *APIHandlers) CreatePlayerHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		PlayerId string `json:"player_id" validate:"required"`
		Region   string `json:"region" validate:"required,len=3"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Fix the request bruh...", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Create Player via Business
	err = a.Logic.CreatePlayer(requestData.PlayerId, requestData.Region)

	if err != nil {
		if err == errormanagement.PlayerIdAlreadyExists {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (a *APIHandlers) CreateRoomHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		PlayerID string `json:"player_id" validate:"required"`
		Mode     string `json:"mode" validate:"required"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Fix the request bruh...", http.StatusBadRequest)
		return
	}

	// Create Room via Business
	room_id, err := a.Logic.CreateRoom(requestData.PlayerID, strings.ToLower(requestData.Mode))
	if err != nil {
		if err == errormanagement.InvalidMode ||
			err == errormanagement.PlayerOccupied {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(room_id))
}

func (a *APIHandlers) GetRoomsHandler(w http.ResponseWriter, r *http.Request) {
	mode := strings.ToLower(r.URL.Query().Get("mode"))
	if mode == "" {
		http.Error(w, "At least type something...", http.StatusBadRequest)
		return
	}
	// Get list of Room Ids via Business
	rooms, err := a.Logic.GetRoomsByMode(mode)
	if err != nil {
		if err == errormanagement.InvalidMode {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	jsonData, _ := json.Marshal(rooms)
	w.Write(jsonData)
}

func (a *APIHandlers) JoinRoomHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		PlayerID string `json:"player_id" validate:"required"`
		RoomID   string `json:"room_id" validate:"required,len=7"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Fix the request bruh...", http.StatusBadRequest)
		return
	}

	// Add Player to room via Business
	err = a.Logic.JoinRoom(requestData.PlayerID, requestData.RoomID)

	if err != nil {
		if err == errormanagement.PlayerNotFound ||
			err == errormanagement.RoomNotFound ||
			err == errormanagement.RoomIsFull ||
			err == errormanagement.PlayerOccupied {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (a *APIHandlers) LeaveRoomHandler(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		PlayerId string `json:"player_id" validate:"required"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, "Fix the request bruh...", http.StatusBadRequest)
		return
	}

	validate := validator.New()
	if err := validate.Struct(requestData); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Remove player from Room via Business
	err = a.Logic.LeaveRoom(r.Context(), requestData.PlayerId)

	if err != nil {
		if err == errormanagement.PlayerIdle ||
			err == errormanagement.PlayerNotFound {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a *APIHandlers) GetModeTrendsByRegion(w http.ResponseWriter, r *http.Request) {
	region := r.URL.Query().Get("region")
	if region == "" {
		http.Error(w, "At least type something...", http.StatusBadRequest)
		return
	}

	// Get List of Mode via Business
	modes, _ := a.Logic.GetModeTrendsByRegion(region)
	jsonData, _ := json.Marshal(modes)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func (a *APIHandlers) GetModeTrendsByRegionV2(w http.ResponseWriter, r *http.Request) {
	var requestData struct {
		PlayerID string `json:"player_id"`
	}

	err := json.NewDecoder(r.Body).Decode(&requestData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Get List of Mode via Business
	modes, err := a.Logic.GetModeTrendsByPlayerRegion(requestData.PlayerID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	jsonData, _ := json.Marshal(modes)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}
