-- name: CreateAuditLog :exec
INSERT INTO audit_logs (
    user_id, action, resource, resource_id,
    method, endpoint, ip_address,
    success, request_body, changes
) VALUES (
    $1, $2, $3, $4,
    $5, $6, $7,
    $8, $9, $10
);

-- name: GetAuditLogsByUser :many
SELECT * FROM audit_logs
WHERE user_id = $1
ORDER BY occurred_at DESC
LIMIT $2 OFFSET $3;

-- name: GetAuditLogsForResource :many
SELECT * FROM audit_logs
WHERE resource = $1 AND resource_id = $2
ORDER BY occurred_at DESC
LIMIT $3 OFFSET $4;
