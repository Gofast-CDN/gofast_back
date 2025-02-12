package helpers

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TestModel implements mgm.Model interface for testing
type TestModel struct {
	// DefaultModel adds _id, created_at and updated_at fields to the Model
	mgm.DefaultModel `bson:",inline"`
}

func (m *TestModel) GetID() interface{} {
	return m.ID
}

func SetupTestDB(t *testing.T) {
	// Use test MongoDB URI or fallback to local dev MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		mongoURI = "mongodb://mongodb:27017"
	}

	err := mgm.SetDefaultConfig(nil, "test_db", options.Client().ApplyURI(mongoURI))
	if err != nil {
		t.Fatalf("Failed to setup MGM: %v", err)
	}

	// Clean database before tests
	cleanDB(t)
}

func cleanDB(t *testing.T) {
	// Use TestModel instead of empty struct
	coll := mgm.Coll(&TestModel{})
	if coll == nil {
		t.Fatal("Database collection not initialized")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := coll.Database().Drop(ctx)
	if err != nil {
		t.Fatalf("Failed to clean test database: %v", err)
	}
}
