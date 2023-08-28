package storage

import (
	"DeathfireArsenal/internal/errormanagement"
	"DeathfireArsenal/pkg/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"math/rand"
	"time"
)

func (s *MongoDBStorage) GetPlayerByID(playerID string) (*models.Player, error) {
	filter := bson.M{"id": playerID}
	var player models.Player
	err := s.playerCollection.FindOne(context.Background(), filter).Decode(&player)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errormanagement.PlayerNotFound
		}
	}
	return &player, nil
}

func (s *MongoDBStorage) PlayerIsAlreadyRegistered(playerID string) bool {
	filter := bson.M{"id": playerID}
	var player models.Player
	err := s.playerCollection.FindOne(context.Background(), filter).Decode(&player)
	if err == mongo.ErrNoDocuments {
		return false
	}
	return true
}

func (s *MongoDBStorage) GetRoomByID(roomID string) (*models.Room, error) {
	filter := bson.M{"id": roomID}

	var room models.Room
	err := s.roomCollection.FindOne(context.Background(), filter).Decode(&room)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errormanagement.RoomNotFound
		}
		return nil, err
	}

	return &room, nil
}

func (s *MongoDBStorage) DeleteRoom(roomID string) error {
	filter := bson.M{"id": roomID}
	_, err := s.roomCollection.DeleteOne(context.Background(), filter)
	return err
}

// Helper function to generate a random room ID.
func generateRoomID() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	b := make([]byte, 7)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}
