package config

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	Client *mongo.Client
	db     *mongo.Database
)

// NewMongoClient creates and returns a new MongoDB client
func NewMongoClient() *mongo.Client {
	uri := os.Getenv("MONGO_URI") // e.g. mongodb://user:pass@localhost:27017
	if uri == "" {
		slog.Error("mongo_uri_not_configured")
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		slog.Error("mongo_connect_failed", "error", err)
		os.Exit(1)
	}

	// Test connection
	if err := client.Ping(ctx, nil); err != nil {
		slog.Error("mongo_ping_failed", "error", err)
		os.Exit(1)
	}

	slog.Info("mongo_connected")

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		slog.Error("mongo_db_name_not_configured")
		os.Exit(1)
	}

	Client = client
	db = client.Database(dbName)

	// Verify required collections exist
	if err := verifyCollections(ctx, db); err != nil {
		slog.Error("mongo_collection_verification_failed", "error", err)
		os.Exit(1)
	}

	slog.Info("mongo_database_selected", "db_name", dbName)
	return client
}

func verifyCollections(ctx context.Context, db *mongo.Database) error {
	collections, err := db.ListCollectionNames(ctx, struct{}{})
	if err != nil {
		return fmt.Errorf("failed to list collections: %w", err)
	}

	required := []string{
		UserColName,
		PostColName,
		CommunityColName,
		CommunityBanColName,
		CommentColName,
		ChannelColName,
		MessageColName,
		VoteColName,
		PollVoteColName,
		NotificationColName,
		ReportColName,
		MembershipColName,
		LikedPostColName,
		SavedPostColName,
		UserPostHistoryColName,
		EmailVerificationColName,
	}

	existing := make(map[string]bool, len(collections))
	for _, c := range collections {
		existing[c] = true
	}

	for _, name := range required {
		if !existing[name] {
			return fmt.Errorf("required collection %q does not exist in database", name)
		}
	}

	slog.Info("mongo_required_collections_verified", "count", len(required))
	return nil
}
