package logic

import (
	"DeathfireArsenal/internal/constants"
	"DeathfireArsenal/internal/errormanagement"
	"DeathfireArsenal/pkg/cache"
	"DeathfireArsenal/pkg/storage"
	"context"
	"errors"
	"fmt"
	"strings"
	"time"
)

type BusinessLogic struct {
	storage *storage.MongoDBStorage
	cache   *cache.RedisCache
}

func NewBusinessLogic(storage *storage.MongoDBStorage, cache *cache.RedisCache) *BusinessLogic {
	return &BusinessLogic{
		storage: storage,
		cache:   cache,
	}
}

func (b *BusinessLogic) CreatePlayer(playerID string, regionCode string) error {
	// Check if Player already exists
	if b.storage.PlayerIsAlreadyRegistered(playerID) {
		return errormanagement.PlayerIdAlreadyExists
	}
	return b.storage.CreatePlayer(playerID, regionCode)
}

func (b *BusinessLogic) CreateRoom(playerID string, mode string) (string, error) {
	//	Check if player exists
	player, err := b.storage.GetPlayerByID(playerID)
	if err != nil {
		return "", err
	}
	//	Check if mode is valid
	if !isValidMode(mode) {
		return "", errormanagement.InvalidMode
	}
	//	Check if player is already in some room
	if len(player.Room) != 0 {
		return "", errormanagement.PlayerOccupied
	}

	room, err := b.storage.CreateRoom(playerID, mode)
	if err != nil {
		return "", err
	}

	b.cache.Invalidate(context.Background())
	return room, nil
}

func (b *BusinessLogic) GetRoomsByMode(mode string) ([]string, error) {
	//	Check if mode is correct
	if !isValidMode(mode) {
		return nil, errormanagement.InvalidMode
	}

	cacheKey := fmt.Sprintf("GetRoomsByMode:%s", mode)

	// Try to get the data from the cache
	var roomIDs []string
	err := b.cache.Get(context.Background(), cacheKey, &roomIDs)
	if err == nil {
		return roomIDs, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		return nil, err
	}

	rooms, err := b.storage.GetRoomsByMode(mode)

	if err != nil {
		return nil, err
	}

	var room_ids []string
	room_ids = make([]string, len(rooms))

	for i := 0; i < len(rooms); i++ {
		room_ids[i] = rooms[i].Id
	}
	b.cache.Set(context.Background(), cacheKey, room_ids, 5*time.Minute)

	return room_ids, nil
}

func (b *BusinessLogic) JoinRoom(playerID string, roomID string) error {
	//	Check if player exists
	player, err := b.storage.GetPlayerByID(playerID)
	if err != nil {
		return err
	}
	//	Check if room exists
	room, err := b.storage.GetRoomByID(roomID)
	if err != nil {
		return err
	}

	//	Check if room is full
	if len(room.PlayerIds) == constants.RoomLimit(constants.ParseMode(room.Mode)) {
		return errormanagement.RoomIsFull
	}

	//	Check if player is already in a room
	if len(player.Room) != 0 {
		return errormanagement.PlayerOccupied
	}

	b.cache.Invalidate(context.Background())
	return b.storage.AddPlayerToRoom(playerID, roomID)
}

func (b *BusinessLogic) LeaveRoom(ctx context.Context, playerID string) error {
	//	Check if player exists
	player, err := b.storage.GetPlayerByID(playerID)
	if err != nil {
		return err
	}

	if len(player.Room) == 0 {
		return errormanagement.PlayerIdle
	}

	b.cache.Invalidate(ctx)
	return b.storage.RemovePlayerFromRoom(ctx, playerID)
}

func (b *BusinessLogic) GetModeTrendsByRegion(region string) (map[string]int, error) {
	cacheKey := fmt.Sprintf("GetModesTrendByRegion:%s", region)

	var modes map[string]int
	err := b.cache.Get(context.Background(), cacheKey, &modes)
	if err == nil {
		return modes, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		return nil, err
	}

	modes, err = b.storage.GetModesByRegionTrend(region)

	if err == nil {
		b.cache.Set(context.Background(), cacheKey, modes, 5*time.Minute)
	}
	return modes, err
}

func (b *BusinessLogic) GetModeTrendsByPlayerRegion(playerId string) (map[string]int, error) {
	cacheKey := fmt.Sprintf("GetModesTrendByPlayerRegion:%s", playerId)

	var modes map[string]int
	err := b.cache.Get(context.Background(), cacheKey, &modes)
	if err == nil {
		return modes, nil
	} else if !errors.Is(err, cache.ErrCacheMiss) {
		return nil, err
	}

	//	Check if player exists
	player, err := b.storage.GetPlayerByID(playerId)
	if err != nil {
		return nil, err
	}
	modes, err = b.storage.GetModesByRegionTrend(player.Region)
	if err == nil {
		b.cache.Set(context.Background(), cacheKey, modes, 5*time.Minute)
	}
	return modes, err
}

// Helper function to check if the mode is valid.
func isValidMode(mode string) bool {
	switch strings.ToLower(mode) {
	case "team deathmatch", "battle royale", "gunsmith", "1 v 1", "mayhem":
		return true
	default:
		return false
	}
}
