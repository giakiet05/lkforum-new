package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type seedUser struct {
	Username string
	Email    string
	Password string
}

var users = []seedUser{
	{
		Username: "e2ee_alice",
		Email:    "e2ee.alice@lkforum.local",
		Password: "E2eeAlice@123",
	},
	{
		Username: "e2ee_bob",
		Email:    "e2ee.bob@lkforum.local",
		Password: "E2eeBob@123",
	},
}

func main() {
	config.LoadConfig()

	client := config.NewMongoClient()
	db := client.Database(config.Cfg.DBName)
	usersCollection := db.Collection(config.UserColName)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	for _, user := range users {
		if err := upsertUser(ctx, usersCollection, user); err != nil {
			log.Fatalf("failed to seed %s: %v", user.Username, err)
		}
	}

	fmt.Println("E2EE test users are ready:")
	for _, user := range users {
		fmt.Printf("username=%s email=%s password=%s\n", user.Username, user.Email, user.Password)
	}
}

func upsertUser(ctx context.Context, collection *mongo.Collection, seed seedUser) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(seed.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	now := time.Now()
	user := model.User{
		ID:         primitive.NewObjectID(),
		Username:   seed.Username,
		Email:      seed.Email,
		Password:   string(hashedPassword),
		Provider:   model.ProviderLocal,
		Role:       model.UserRole,
		Reputation: 0,
		Settings:   model.NewDefaultSettings(),
		IsVerified: true,
		IsBanned:   false,
		CreatedAt:  now,
		RoleContent: model.RoleContent{
			AsUser: &model.UserRoleContent{
				Avatar: &model.Image{},
				Stats: &model.ActivityStats{
					PostCount:    0,
					CommentCount: 0,
					TotalUpvotes: 0,
					JoinedAt:     now,
					LastActiveAt: now,
				},
			},
		},
	}

	filter := bson.M{
		"$or": []bson.M{
			{"email": seed.Email},
			{"username": seed.Username},
		},
	}
	update := bson.M{
		"$set": bson.M{
			"username":    user.Username,
			"email":       user.Email,
			"password":    user.Password,
			"provider":    user.Provider,
			"role":        user.Role,
			"reputation":  user.Reputation,
			"settings":    user.Settings,
			"is_verified": user.IsVerified,
			"is_banned":   user.IsBanned,
		},
		"$setOnInsert": bson.M{
			"_id":          user.ID,
			"created_at":   user.CreatedAt,
			"role_content": user.RoleContent,
		},
		"$unset": bson.M{
			"deleted_at": "",
			"ban_until":  "",
			"ban_reason": "",
		},
	}
	opts := options.Update().SetUpsert(true)

	if _, err := collection.UpdateOne(ctx, filter, update, opts); err != nil {
		return fmt.Errorf("upsert user: %w", err)
	}

	return nil
}
