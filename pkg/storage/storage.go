package storage

import (
	"DeathfireArsenal/pkg/models"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"sort"
)

type MongoDBStorage struct {
	roomCollection   *mongo.Collection
	playerCollection *mongo.Collection
}

func NewMongoDBStorage(rooms *mongo.Collection, players *mongo.Collection) *MongoDBStorage {
	return &MongoDBStorage{
		roomCollection:   rooms,
		playerCollection: players,
	}
}

func (s *MongoDBStorage) CreatePlayer(playerId string, regionCode string) error {
	player := models.Player{Id: playerId, Region: regionCode, Room: ""}
	_, err := s.playerCollection.InsertOne(context.Background(), player)
	return err
}

func (s *MongoDBStorage) CreateRoom(playerId string, mode string) (string, error) {
	filter := bson.M{"id": playerId}
	playerList := []string{playerId}

	//	Create room
	random_room_id := generateRoomID()
	room := models.Room{Id: random_room_id, Mode: mode, PlayerIds: playerList}
	update := bson.M{"$set": bson.M{"room": random_room_id}}
	_, err := s.playerCollection.UpdateOne(context.Background(), filter, update)

	if err != nil {
		return "", err
	}

	_, err = s.roomCollection.InsertOne(context.Background(), room)
	if err != nil {

		// Reverting previous changes as well - That's attention to detail!
		update := bson.M{"$set": bson.M{"room": ""}}
		s.playerCollection.UpdateOne(context.Background(), filter, update)

		return "", err
	}
	return random_room_id, err
}

func (s *MongoDBStorage) GetRoomsByMode(mode string) ([]*models.Room, error) {
	filter := bson.M{"mode": mode}

	cur, err := s.roomCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	var rooms []*models.Room
	for cur.Next(context.Background()) {
		var room models.Room
		err := cur.Decode(&room)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return rooms, nil
}

func (s *MongoDBStorage) AddPlayerToRoom(playerId string, roomID string) error {
	playerFilter := bson.M{"id": playerId}
	update := bson.M{"$set": bson.M{"room": roomID}}
	_, err := s.playerCollection.UpdateOne(context.Background(), playerFilter, update)
	if err != nil {
		return err
	}
	roomFilter := bson.M{"id": roomID}
	update = bson.M{"$addToSet": bson.M{"playerids": playerId}}
	_, err = s.roomCollection.UpdateOne(context.Background(), roomFilter, update)

	if err != nil {
		update := bson.M{"$set": bson.M{"room": ""}}
		s.playerCollection.UpdateOne(context.Background(), playerFilter, update)
	}
	return err
}
func (s *MongoDBStorage) RemovePlayerFromRoom(ctx context.Context, playerId string) error {
	// Find the room that the player is currently in.
	playerFilter := bson.M{"id": playerId}
	var player models.Player
	err := s.playerCollection.FindOne(ctx, playerFilter).Decode(&player)
	// Remove the player from the playerIds list of the room.
	roomFilter := bson.M{"id": player.Room}
	update := bson.M{"$pull": bson.M{"playerids": playerId}}

	_, err = s.roomCollection.UpdateOne(ctx, roomFilter, update)
	if err != nil {
		return err
	} else {
		var room models.Room
		err = s.roomCollection.FindOne(ctx, roomFilter).Decode(&room)
		if err != nil {
			return err
		}
		if len(room.PlayerIds) == 0 {
			s.roomCollection.DeleteOne(ctx, roomFilter)
		}
	}

	// Update the player's room field to empty.
	playerUpdate := bson.M{"$set": bson.M{"room": ""}}
	_, err = s.playerCollection.UpdateOne(ctx, playerFilter, playerUpdate)
	return err
}

func (s *MongoDBStorage) GetModesByRegionTrend(region string) (map[string]int, error) {
	// Define the filter to find players with non-empty regions
	filter := bson.M{"room": bson.M{"$ne": ""}, "region": region}
	ans := make(map[string]int)

	// Execute the query and get the cursor
	cursor, err := s.playerCollection.Find(context.Background(), filter)
	if err != nil {
		return nil, err
	}

	for cursor.Next(context.Background()) {
		var player models.Player
		if err := cursor.Decode(&player); err != nil {
			return nil, err
		}

		var room models.Room
		if err := s.roomCollection.FindOne(context.Background(), bson.M{"id": player.Room}).Decode(&room); err != nil {
			return nil, err
		}

		ans[room.Mode]++
	}

	type ModeCount struct {
		Mode  string
		Count int
	}
	var modeCounts []ModeCount
	for mode, count := range ans {
		modeCounts = append(modeCounts, ModeCount{Mode: mode, Count: count})
	}

	sort.Slice(modeCounts, func(i, j int) bool {
		return modeCounts[i].Count > modeCounts[j].Count
	})

	top3Map := make(map[string]int)
	for i, entry := range modeCounts {
		top3Map[entry.Mode] = entry.Count
		if i == 2 {
			break
		}
	}

	return top3Map, nil
}
