package service

import (
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
		_ = s.repo.Create(ctx, entry)
	}()
}
