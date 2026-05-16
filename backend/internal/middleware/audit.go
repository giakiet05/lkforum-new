package middleware

import (
	"github.com/giakiet05/lkforum/internal/auth"
	"github.com/giakiet05/lkforum/internal/model"
	"github.com/giakiet05/lkforum/internal/service"
	"github.com/gin-gonic/gin"
)

var auditLogService service.AuditLogService

func SetAuditLogService(s service.AuditLogService) {
	auditLogService = s
}

func RecordAudit(c *gin.Context, action, targetType, targetID, reason string, metadata map[string]interface{}) {
	if auditLogService == nil {
		return
	}

	entry := &model.AuditLog{
		Action:     action,
		TargetType: targetType,
		TargetID:   targetID,
		Reason:     reason,
		IP:         c.ClientIP(),
		UserAgent:  c.Request.UserAgent(),
		Metadata:   metadata,
	}
	if val, ok := c.Get("authUser"); ok {
		if user, ok := val.(auth.AuthUser); ok {
			entry.ActorID = user.ID
			entry.ActorRole = string(user.Role)
		}
	}

	auditLogService.Record(entry)
}
