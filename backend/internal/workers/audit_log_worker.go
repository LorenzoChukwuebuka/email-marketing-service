package worker

import (
	"context"
	"database/sql"
	db "email-marketing-service/internal/db/sqlc"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/sqlc-dev/pqtype"
	"log"
	"net"
	"net/http"
)

type AuditAction string

const (
	AuditActionCreate      AuditAction = "CREATE"
	AuditActionUpdate      AuditAction = "UPDATE"
	AuditActionDelete      AuditAction = "DELETE"
	AuditActionLogin       AuditAction = "LOGIN"
	AuditActionLogout      AuditAction = "LOGOUT"
	AuditActionLoginFailed AuditAction = "LOGIN_FAILED"
)

func (w *Worker) ProcessLogAuditCreate(ctx context.Context, payload AuditCreatePayload) error {
	reqJSON, _ := json.Marshal(payload.RequestBody)

	var resourceUUID uuid.NullUUID
	if payload.ResourceID != nil {
		resourceUUID = uuid.NullUUID{UUID: *payload.ResourceID, Valid: true}
	}

	err := w.Store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		UserID:      payload.UserID,
		Action:      AuditActionCreate,
		Resource:    payload.Resource,
		ResourceID:  resourceUUID,
		Method:      sql.NullString{String: payload.Method, Valid: payload.Method != ""},
		Endpoint:    sql.NullString{String: payload.Endpoint, Valid: payload.Endpoint != ""},
		IpAddress:   ipToInet(payload.IP),
		Success:     sql.NullBool{Bool: payload.Success, Valid: true},
		RequestBody: pqtype.NullRawMessage{RawMessage: reqJSON, Valid: len(reqJSON) > 0},
		Changes:     pqtype.NullRawMessage{Valid: false},
	})

	if err != nil {
		log.Printf("Failed to log audit create: %v", err)
		return err
	}

	log.Printf("Audit log created for resource: %s", payload.Resource)
	return nil
}

func (w *Worker) ProcessLogAuditUpdate(ctx context.Context, payload AuditUpdatePayload) error {
	reqJSON, _ := json.Marshal(payload.RequestBody)
	changesJSON, _ := json.Marshal(payload.Changes)

	var resourceUUID uuid.NullUUID
	if payload.ResourceID != nil {
		resourceUUID = uuid.NullUUID{UUID: *payload.ResourceID, Valid: true}
	}

	err := w.Store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		UserID:      payload.UserID,
		Action:      AuditActionUpdate,
		Resource:    payload.Resource,
		ResourceID:  resourceUUID,
		Method:      sql.NullString{String: payload.Method, Valid: payload.Method != ""},
		Endpoint:    sql.NullString{String: payload.Endpoint, Valid: payload.Endpoint != ""},
		IpAddress:   ipToInet(payload.IP),
		Success:     sql.NullBool{Bool: payload.Success, Valid: true},
		RequestBody: pqtype.NullRawMessage{RawMessage: reqJSON, Valid: len(reqJSON) > 0},
		Changes:     pqtype.NullRawMessage{RawMessage: changesJSON, Valid: len(changesJSON) > 0},
	})

	if err != nil {
		log.Printf("Failed to log audit update: %v", err)
		return err
	}

	log.Printf("Audit log updated for resource: %s", payload.Resource)
	return nil
}

func (w *Worker) ProcessLogAuditDelete(ctx context.Context, payload AuditDeletePayload) error {
	reqJSON, _ := json.Marshal(payload.RequestBody)

	var resourceUUID uuid.NullUUID
	if payload.ResourceID != nil {
		resourceUUID = uuid.NullUUID{UUID: *payload.ResourceID, Valid: true}
	}

	err := w.Store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		UserID:      payload.UserID,
		Action:      AuditActionDelete,
		Resource:    payload.Resource,
		ResourceID:  resourceUUID,
		Method:      sql.NullString{String: payload.Method, Valid: payload.Method != ""},
		Endpoint:    sql.NullString{String: payload.Endpoint, Valid: payload.Endpoint != ""},
		IpAddress:   ipToInet(payload.IP),
		Success:     sql.NullBool{Bool: payload.Success, Valid: true},
		RequestBody: pqtype.NullRawMessage{RawMessage: reqJSON, Valid: len(reqJSON) > 0},
		Changes:     pqtype.NullRawMessage{Valid: false},
	})

	if err != nil {
		log.Printf("Failed to log audit delete: %v", err)
		return err
	}

	log.Printf("Audit log deleted for resource: %s", payload.Resource)
	return nil
}

func (w *Worker) ProcessLogAuditLogin(ctx context.Context, payload AuditLoginPayload) error {
	body := map[string]string{"username": payload.Username}
	reqJSON, _ := json.Marshal(body)

	err := w.Store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		UserID:      payload.UserID,
		Action:      AuditActionLogin,
		Resource:    "UserAuth",
		ResourceID:  uuid.NullUUID{UUID: payload.UserID, Valid: true},
		Method:      sql.NullString{String: payload.Method, Valid: payload.Method != ""},
		Endpoint:    sql.NullString{String: payload.Endpoint, Valid: payload.Endpoint != ""},
		IpAddress:   ipToInet(payload.IP),
		Success:     sql.NullBool{Bool: true, Valid: true},
		RequestBody: pqtype.NullRawMessage{RawMessage: reqJSON, Valid: len(reqJSON) > 0},
		Changes:     pqtype.NullRawMessage{Valid: false},
	})

	if err != nil {
		log.Printf("Failed to log audit login: %v", err)
		return err
	}

	log.Printf("Audit login logged for user: %s", payload.UserID)
	return nil
}

func (w *Worker) ProcessLogAuditFailedLogin(ctx context.Context, payload AuditFailedLoginPayload) error {
	body := map[string]string{"username": payload.AttemptedUsername}
	reqJSON, _ := json.Marshal(body)

	err := w.Store.CreateAuditLog(ctx, db.CreateAuditLogParams{
		UserID:      uuid.Nil,
		Action:      AuditActionLoginFailed,
		Resource:    "UserAuth",
		ResourceID:  uuid.NullUUID{Valid: false},
		Method:      sql.NullString{String: payload.Method, Valid: payload.Method != ""},
		Endpoint:    sql.NullString{String: payload.Endpoint, Valid: payload.Endpoint != ""},
		IpAddress:   ipToInet(payload.IP),
		Success:     sql.NullBool{Bool: false, Valid: true},
		RequestBody: pqtype.NullRawMessage{RawMessage: reqJSON, Valid: len(reqJSON) > 0},
		Changes:     pqtype.NullRawMessage{Valid: false},
	})

	if err != nil {
		log.Printf("Failed to log audit failed login: %v", err)
		return err
	}

	log.Printf("Audit failed login logged for username: %s", payload.AttemptedUsername)
	return nil
}

// ============================================
// HELPER FUNCTIONS
// ============================================

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

func (w *Worker) GetClientIP(r *http.Request) net.IP {
	ipStr := r.Header.Get("X-Forwarded-For")
	if ipStr == "" {
		ipStr, _, _ = net.SplitHostPort(r.RemoteAddr)
	}
	return net.ParseIP(ipStr)
}
