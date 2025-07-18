-- name: CreateSupportTicket :one
INSERT INTO
    support_tickets (
        user_id,
        name,
        email,
        subject,
        description,
        ticket_number,
        status,
        priority
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8
    ) RETURNING *;

-- name: CreateTicketMessage :one
INSERT INTO
    ticket_messages (
        ticket_id,
        user_id,
        message,
        is_admin
    )
VALUES ($1, $2, $3, $4) RETURNING *;

-- name: CreateTicketFile :one
INSERT INTO
    ticket_files (
        message_id,
        file_name,
        file_path
    )
VALUES ($1, $2, $3) RETURNING *;

-- name: FindTicketByID :one
SELECT * FROM support_tickets WHERE id = $1;

-- name: UpdateTicketStatus :one
UPDATE support_tickets
SET
    status = $2,
    last_reply = $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING *;

-- name: GetTicketWithMessages :many
SELECT
    t.id as ticket_id,
    t.user_id as ticket_user_id,
    t.name as ticket_name,
    t.email as ticket_email,
    t.subject as ticket_subject,
    t.description as ticket_description,
    t.ticket_number,
    t.status as ticket_status,
    t.priority as ticket_priority,
    t.last_reply as ticket_last_reply,
    t.created_at as ticket_created_at,
    t.updated_at as ticket_updated_at,
    m.id as message_id,
    m.user_id as message_user_id,
    m.message,
    m.is_admin,
    m.created_at as message_created_at,
    m.updated_at as message_updated_at
FROM
    support_tickets t
    LEFT JOIN ticket_messages m ON t.id = m.ticket_id
WHERE
    t.id = $1
ORDER BY m.created_at ASC;

-- name: GetTicketsByUserID :many
SELECT *
FROM support_tickets
WHERE
    user_id = $1
ORDER BY created_at DESC;

-- name: GetMessageFiles :many
SELECT * FROM ticket_files WHERE message_id = $1;

-- name: GetTicketFiles :many
SELECT tf.*
FROM
    ticket_files tf
    JOIN ticket_messages tm ON tf.message_id = tm.id
WHERE
    tm.ticket_id = $1;

-- name: CloseTicketByID :one
UPDATE support_tickets
SET
    status = 'closed',
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = $1 RETURNING *;

-- name: CloseStaleTickets :many
UPDATE support_tickets
SET
    status = 'closed',
    updated_at = CURRENT_TIMESTAMP
WHERE
    status NOT IN('closed', 'resolved')
    AND (
        -- Tickets with no replies that are older than 48 hours
        (
            last_reply IS NULL
            AND created_at < CURRENT_TIMESTAMP - INTERVAL '48 hours'
        )
        OR
        -- Tickets with replies but last reply is older than 48 hours
        (
            last_reply IS NOT NULL
            AND last_reply < CURRENT_TIMESTAMP - INTERVAL '48 hours'
        )
    ) RETURNING *;

-- name: GetAllTicketsWithPagination :many
SELECT 
   *
FROM support_tickets
WHERE 
    CASE 
        WHEN $1::text != '' THEN 
            (subject ILIKE '%' || $1 || '%' OR 
             ticket_number ILIKE '%' || $1 || '%' OR 
             name ILIKE '%' || $1 || '%' OR 
             email ILIKE '%' || $1 || '%')
        ELSE TRUE
    END
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetPendingTicketsWithPagination :many
SELECT 
     *
FROM support_tickets
WHERE 
    status = 'pending'
    AND CASE 
        WHEN $1::text != '' THEN 
            (subject ILIKE '%' || $1 || '%' OR 
             ticket_number ILIKE '%' || $1 || '%' OR 
             name ILIKE '%' || $1 || '%' OR 
             email ILIKE '%' || $1 || '%')
        ELSE TRUE
    END
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: GetClosedTicketsWithPagination :many
SELECT 
    *
FROM support_tickets
WHERE 
    status = 'closed'
    AND CASE 
        WHEN $1::text != '' THEN 
            (subject ILIKE '%' || $1 || '%' OR 
             ticket_number ILIKE '%' || $1 || '%' OR 
             name ILIKE '%' || $1 || '%' OR 
             email ILIKE '%' || $1 || '%')
        ELSE TRUE
    END
ORDER BY created_at DESC
LIMIT $2 OFFSET $3;


-- name: GetAllTicketsCount :one
SELECT COUNT(*) FROM support_tickets;

-- name: GetPendingTicketsCount :one
SELECT COUNT(*) FROM support_tickets WHERE status = 'pending';

-- name: GetClosedTicketsCount :one
SELECT COUNT(*) FROM support_tickets WHERE status = 'closed';