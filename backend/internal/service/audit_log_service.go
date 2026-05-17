package service

import (
	"log/slog"
	"time"

	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/repo"
	"github.com/giakiet05/lkforum/internal/util"
)

type AuditLogService interface {
	Record(entry *model.AuditLog)
}

type auditLogService struct {
	repo repo.AuditLogRepo
}

func NewAuditLogService(repo repo.AuditLogRepo) AuditLogService {
	return &auditLogService{repo: repo}
}

func (s *auditLogService) Record(entry *model.AuditLog) {
	if entry == nil {
		return
	}
	if entry.CreatedAt.IsZero() {
		entry.CreatedAt = time.Now()
	}

	go func() {
		ctx, cancel := util.NewDefaultDBContext()
		defer cancel()
		if err := s.repo.Create(ctx, entry); err != nil {
			slog.ErrorContext(
				ctx,
				"audit_log_persist_failed",
				"action", entry.Action,
				"actor_id", entry.ActorID,
				"target_type", entry.TargetType,
				"target_id", entry.TargetID,
				"error", err,
			)
			return
		}

		slog.InfoContext(
			ctx,
			"audit_log_persisted",
			"action", entry.Action,
			"actor_id", entry.ActorID,
			"target_type", entry.TargetType,
			"target_id", entry.TargetID,
		)
	}()
}
