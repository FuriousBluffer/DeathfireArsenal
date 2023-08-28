package errormanagement

import "errors"

var (
	PlayerIdAlreadyExists = errors.New("Player ID is already taken")
	PlayerNotFound        = errors.New("Player does not exist")
	RoomNotFound          = errors.New("Room not found")
	RoomIsFull            = errors.New("Current Room is full")
	PlayerOccupied        = errors.New("Player is already in a combat for his virtual life. Can't afford to join another.")
	PlayerIdle            = errors.New("Player not playing any game rn...")
	InvalidMode           = errors.New("This Mode of Game does not exist.")
)
