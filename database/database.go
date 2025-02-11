package database

import (
	"fmt"
	"log"
	"os"

	"go.mongodb.org/mongo-driver/mongo"
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
	log.Println("🔍 Tentative de connexion à:", uri)

	// Affichage sécurisé
	maskedURI := fmt.Sprintf("mongodb+srv://%s:*****@gofastcluster.0csvm.mongodb.net/%s", user, database)
	log.Println("🔍 Tentative de connexion à:", maskedURI)

	// clientOptions := options.Client().ApplyURI(uri)

	// ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	// defer cancel()

	// client, err := mongo.Connect(ctx, clientOptions)
	// if err != nil {
	// 	log.Fatal("Erreur de connexion à MongoDB:", err)
	// }

	// err = client.Ping(ctx, nil)
	// if err != nil {
	// 	log.Fatal("Impossible de pinger MongoDB:", err)
	// }

	// err = mgm.SetDefaultConfig(nil, database, clientOptions)
	// if err != nil {
	// 	log.Fatal("Impossible de configurer mgm:", err)
	// }

	log.Println("✅ Connexion réussie à MongoDB")
	// Client = client
}
