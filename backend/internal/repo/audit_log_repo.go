package repo

import (
	"context"

	"github.com/giakiet05/lkforum/internal/config"
	"github.com/giakiet05/lkforum/internal/model"
	"go.mongodb.org/mongo-driver/mongo"
)

type AuditLogRepo interface {
	Create(ctx context.Context, log *model.AuditLog) error
}

type auditLogRepo struct {
	collection *mongo.Collection
}

func NewAuditLogRepo(db *mongo.Database) AuditLogRepo {
	return &auditLogRepo{collection: db.Collection(config.AuditLogColName)}
}

func (r *auditLogRepo) Create(ctx context.Context, log *model.AuditLog) error {
	_, err := r.collection.InsertOne(ctx, log)
	return err
}
