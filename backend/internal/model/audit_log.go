package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AuditLog struct {
	ID         primitive.ObjectID     `bson:"_id,omitempty" json:"id"`
	ActorID    string                 `bson:"actor_id,omitempty" json:"actor_id,omitempty"`
	ActorRole  string                 `bson:"actor_role,omitempty" json:"actor_role,omitempty"`
	Action     string                 `bson:"action" json:"action"`
	TargetType string                 `bson:"target_type,omitempty" json:"target_type,omitempty"`
	TargetID   string                 `bson:"target_id,omitempty" json:"target_id,omitempty"`
	Reason     string                 `bson:"reason,omitempty" json:"reason,omitempty"`
	IP         string                 `bson:"ip,omitempty" json:"ip,omitempty"`
	UserAgent  string                 `bson:"user_agent,omitempty" json:"user_agent,omitempty"`
	Metadata   map[string]interface{} `bson:"metadata,omitempty" json:"metadata,omitempty"`
	CreatedAt  time.Time              `bson:"created_at" json:"created_at"`
}
