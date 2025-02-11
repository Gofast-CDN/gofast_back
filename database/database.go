package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var MongoClient *mongo.Client
var RedisClient *redis.Client

func ConnectMongoDB() {
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")
	database := os.Getenv("MONGO_DATABASE")

	uri := fmt.Sprintf("mongodb+srv://%s:%s@gofastcluster.0csvm.mongodb.net/%s?retryWrites=true&w=majority&appName=GoFastCluster",
		user, password, database)

	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Erreur de connexion à MongoDB:", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Impossible de pinger MongoDB:", err)
	}

	log.Println("Connexion réussie à MongoDB")
	MongoClient = client
}

func ConnectRedis() {
	redisURL := os.Getenv("REDIS_URL")
	redisPASSWORD := os.Getenv("REDIS_PASSWORD")

	if redisURL == "" {
		log.Fatal("REDIS_URL not defined on environment variables")
	}

	RedisClient = redis.NewClient(&redis.Options{
		Addr:     redisURL,
		Password: redisPASSWORD,
		DB:       0,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Impossible de se connecter à Redis:", err)
	}

	log.Println("Connexion réussie à Redis")
}

func Connect() {
	ConnectMongoDB()
	ConnectRedis()
}
