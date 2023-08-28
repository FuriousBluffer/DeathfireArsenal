package storage

type Storage interface {
	CreatePlayer(playerId string, region string) error
	CreateRoom(playerId string, mode string) error
	GetRoomsByMode(mode string) error
	AddPlayerToRoom(playerId string, roomID string) error
	RemovePlayerFromRoom(playerId string) error
	GetModesByRegionTrend(region string) error
}
