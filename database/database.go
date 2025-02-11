package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var Client *mongo.Client

func Connect() {
	user := os.Getenv("MONGO_USER")
	password := os.Getenv("MONGO_PASSWORD")
	database := os.Getenv("MONGO_DATABASE")

	// Vérification des variables d'environnement
	if user == "" || password == "" || database == "" {
		log.Fatal("❌ Les variables d'environnement MONGO_USER, MONGO_PASSWORD ou MONGO_DATABASE ne sont pas définies.")
	}

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

	err = mgm.SetDefaultConfig(nil, database, clientOptions)
	if err != nil {
		log.Fatal("Impossible de configurer mgm:", err)
	}

	log.Println("✅ Connexion réussie à MongoDB")
	Client = client
}
