-- name: CreateInvitation :one
INSERT INTO invitations (
    company_id,
    invited_by,
    email,
    token,
    expires_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetInvitationByToken :one
SELECT 
    i.*,
    c.companyname,
    u.fullname as invited_by_name
FROM invitations i
JOIN companies c ON i.company_id = c.id
JOIN users u ON i.invited_by = u.id
WHERE i.token = $1 AND i.deleted_at IS NULL;

-- name: GetInvitationByID :one
SELECT * FROM invitations 
WHERE id = $1 AND deleted_at IS NULL;

-- name: GetInvitationsByCompany :many
SELECT 
    i.*,
    u.fullname as invited_by_name
FROM invitations i
JOIN users u ON i.invited_by = u.id
WHERE i.company_id = $1 AND i.deleted_at IS NULL
ORDER BY i.created_at DESC;

-- name: GetPendingInvitationsByEmail :many
SELECT 
    i.*,
    c.companyname,
    u.fullname as invited_by_name
FROM invitations i
JOIN companies c ON i.company_id = c.id
JOIN users u ON i.invited_by = u.id
WHERE i.email = $1 AND i.status = 'pending' AND i.expires_at > now() AND i.deleted_at IS NULL;

-- name: UpdateInvitationStatus :one
UPDATE invitations 
SET 
    status = $2,
    accepted_at = CASE WHEN $2 = 'accepted' THEN now() ELSE accepted_at END,
    accepted_by = $3,
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL
RETURNING *;

-- name: GetUserWhoAcceptedInvitation :one
SELECT 
    u.*,
    i.accepted_at,
    i.created_at as invitation_created_at
FROM invitations i
JOIN users u ON i.accepted_by = u.id
WHERE i.id = $1 AND i.status = 'accepted' AND i.deleted_at IS NULL;

-- name: GetInvitationByAcceptedUser :one
SELECT 
    i.*,
    invited_by_user.fullname as invited_by_name,
    c.companyname
FROM invitations i
JOIN users invited_by_user ON i.invited_by = invited_by_user.id
JOIN companies c ON i.company_id = c.id
WHERE i.accepted_by = $1 AND i.status = 'accepted' AND i.deleted_at IS NULL;

-- name: GetUsersFromInvitations :many
SELECT 
    u.*,
    i.created_at as invitation_created_at,
    i.accepted_at,
    invited_by_user.fullname as invited_by_name
FROM invitations i
JOIN users u ON i.accepted_by = u.id
JOIN users invited_by_user ON i.invited_by = invited_by_user.id
WHERE i.company_id = $1 AND i.status = 'accepted' AND i.deleted_at IS NULL
ORDER BY i.accepted_at DESC;

-- name: ExpireInvitation :exec
UPDATE invitations 
SET 
    status = 'expired',
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: CancelInvitation :exec
UPDATE invitations 
SET 
    status = 'cancelled',
    updated_at = now()
WHERE id = $1 AND deleted_at IS NULL;

-- name: DeleteInvitation :exec
UPDATE invitations 
SET 
    deleted_at = now(),
    updated_at = now()
WHERE id = $1;

-- name: CleanupExpiredInvitations :exec
UPDATE invitations 
SET 
    status = 'expired',
    updated_at = now()
WHERE status = 'pending' AND expires_at < now() AND deleted_at IS NULL;

-- name: GetInvitationStats :one
SELECT 
    COUNT(*) as total_invitations,
    COUNT(CASE WHEN status = 'pending' THEN 1 END) as pending_invitations,
    COUNT(CASE WHEN status = 'accepted' THEN 1 END) as accepted_invitations,
    COUNT(CASE WHEN status = 'expired' THEN 1 END) as expired_invitations,
    COUNT(CASE WHEN status = 'cancelled' THEN 1 END) as cancelled_invitations
FROM invitations 
WHERE company_id = $1 AND deleted_at IS NULL;
