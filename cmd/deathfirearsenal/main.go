package main

import (
	"DeathfireArsenal/internal/logic"
	api_handlers "DeathfireArsenal/pkg/api"
	"DeathfireArsenal/pkg/cache"
	"DeathfireArsenal/pkg/storage"
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	err := godotenv.Load("./internal/env/.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Mongo Setup
	mongoClient, err := mongo.NewClient(options.Client().ApplyURI(os.Getenv("MONGODB_URL")))
	if err != nil {
		log.Fatal("Failed to create MongoDB client:", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = mongoClient.Connect(ctx)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB:", err)
	}
	defer mongoClient.Disconnect(ctx)
	roomCollection := mongoClient.Database("DeathfireArsenal").Collection("rooms")
	playerCollection := mongoClient.Database("DeathfireArsenal").Collection("players")

	// Redis Setup
	redisClient := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: "",
		DB:       0,
	})
	redisClient.FlushAll(ctx)

	mongoDBStorage := storage.NewMongoDBStorage(roomCollection, playerCollection)
	redisCache := cache.NewRedisCache(redisClient)
	businessLogic := logic.NewBusinessLogic(mongoDBStorage, redisCache)
	apiHandlers := api_handlers.APIHandlers{
		Logic: businessLogic,
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/createPlayer", apiHandlers.CreatePlayerHandler).Methods("POST")
	router.HandleFunc("/api/createRoom", apiHandlers.CreateRoomHandler).Methods("POST")
	router.HandleFunc("/api/getRooms", apiHandlers.GetRoomsHandler).Methods("GET")
	router.HandleFunc("/api/joinRoom", apiHandlers.JoinRoomHandler).Methods("POST")
	router.HandleFunc("/api/leaveRoom", apiHandlers.LeaveRoomHandler).Methods("POST")
	router.HandleFunc("/api/getModeTrendsByRegion", apiHandlers.GetModeTrendsByRegion).Methods("GET")
	router.HandleFunc("/api/getModeTrendsByRegionV2", apiHandlers.GetModeTrendsByRegionV2).Methods("GET")

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt)
	go func() {
		<-stopChan
		log.Println("Shutting down the server...")
		server.Shutdown(context.Background())
	}()

	fmt.Println("Server is now running on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal("Server failed to start:", err)
	}
}
