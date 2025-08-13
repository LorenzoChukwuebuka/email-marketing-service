package logger

import (
	"context"
	"database/sql"
	db "email-marketing-service/internal/db/sqlc"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"net"
	"net/http"
)

type AuditLogger struct {
	q db.Store
}

func NewAuditLogger(q db.Store) *AuditLogger {
	return &AuditLogger{q: q}
}

type Change struct {
	Old any `json:"old"`
	New any `json:"new"`
}

func (l *AuditLogger) logEvent(
	ctx context.Context,
	userID uuid.UUID,
	action db.AuditAction,
	resource string,
	resourceID *uuid.UUID,
	method string,
	endpoint string,
	ip net.IP,
	success bool,
	requestBody interface{},
	changes map[string]Change,
) error {
	reqJSON, _ := json.Marshal(requestBody)
	changesJSON, _ := json.Marshal(changes)

	return l.q.CreateAuditLog(ctx, db.CreateAuditLogParams{
		UserID:      userID,
		Action:      action,
		Resource:    resource,
		ResourceID:  uuid.NullUUID{UUID: uuid.Nil, Valid: resourceID != nil},
		Method:      sql.NullString{String: method, Valid: method != ""},
		Endpoint:    sql.NullString{String: endpoint, Valid: endpoint != ""},
		IpAddress:   ipToInet(ip),
		Success:     sql.NullBool{Bool: success, Valid: true},
		RequestBody: pqtype.NullRawMessage{RawMessage: reqJSON, Valid: len(reqJSON) > 0},
		Changes:     pqtype.NullRawMessage{RawMessage: changesJSON, Valid: len(changesJSON) > 0},
	})
}

// Helper for HTTP handlers
func (l *AuditLogger) GetClientIP(r *http.Request) net.IP {
	ipStr := r.Header.Get("X-Forwarded-For")
	if ipStr == "" {
		ipStr, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return net.ParseIP(ipStr)
}

func ipToInet(ip net.IP) pqtype.Inet {
	var cidr string
	if ip.To4() != nil {
		cidr = ip.String() + "/32"
	} else {
		cidr = ip.String() + "/128"
	}
	_, ipNet, _ := net.ParseCIDR(cidr)
	return pqtype.Inet{
		IPNet: *ipNet,
		Valid: true,
	}
}

// CREATE
func (l *AuditLogger) LogCreate(ctx context.Context, userID uuid.UUID, resource string, resourceID *uuid.UUID, method, endpoint string, ip net.IP, requestBody interface{}) error {
	return l.logEvent(ctx, userID, db.AuditActionCREATE, resource, resourceID, method, endpoint, ip, true, requestBody, nil)
}

// UPDATE
func (l *AuditLogger) LogUpdate(ctx context.Context, userID uuid.UUID, resource string, resourceID *uuid.UUID, method, endpoint string, ip net.IP, requestBody interface{}, changes map[string]Change) error {
	return l.logEvent(ctx, userID, db.AuditActionUPDATE, resource, resourceID, method, endpoint, ip, true, requestBody, changes)
}

// DELETE
func (l *AuditLogger) LogDelete(ctx context.Context, userID uuid.UUID, resource string, resourceID *uuid.UUID, method, endpoint string, ip net.IP, requestBody interface{}) error {
	return l.logEvent(ctx, userID, db.AuditActionDELETE, resource, resourceID, method, endpoint, ip, true, requestBody, nil)
}

// LOGIN success
func (l *AuditLogger) LogLogin(ctx context.Context, userID uuid.UUID, method, endpoint string, ip net.IP, username string) error {
	body := map[string]string{"username": username}
	return l.logEvent(ctx, userID, db.AuditActionLOGIN, "UserAuth", &userID, method, endpoint, ip, true, body, nil)
}

// LOGIN failed
func (l *AuditLogger) LogFailedLogin(ctx context.Context, attemptedUsername string, method, endpoint string, ip net.IP) error {
	body := map[string]string{"username": attemptedUsername}
	return l.logEvent(ctx, uuid.Nil, db.AuditActionLOGINFAILED, "UserAuth", nil, method, endpoint, ip, false, body, nil)
}
