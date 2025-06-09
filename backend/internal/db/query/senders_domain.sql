-- name: CreateDomain :one
INSERT INTO
    domains (
        user_id,
        company_id,
        domain,
        txt_record,
        dmarc_record,
        dkim_selector,
        dkim_public_key,
        dkim_private_key,
        spf_record,
        verified,
        mx_record
    )
VALUES (
        $1,
        $2,
        $3,
        $4,
        $5,
        $6,
        $7,
        $8,
        $9,
        $10,
        $11
    ) RETURNING *;

-- name: CreateSender :one
INSERT INTO
    senders (
        user_id,
        company_id,
        name,
        email,
        verified,
        is_signed,
        domain_id
    )
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: CheckSenderExists :one
SELECT EXISTS (
        SELECT 1
        FROM senders
        WHERE
            email = $1
            AND name = $2
            AND company_id = $3
            AND deleted_at IS NULL
    ) AS exists;

-- name: CheckDomainExists :one
SELECT EXISTS (
        SELECT 1
        FROM domains
        WHERE
            domain = $1
            and company_id = $2
            AND deleted_at is NULL
    ) AS exists;

-- name: UpdateDomain :exec
UPDATE domains
SET
    domain = COALESCE(@domain, domain),
    txt_record = COALESCE(@txt_record, txt_record),
    dmarc_record = COALESCE(@dmarc_record, dmarc_record),
    dkim_selector = COALESCE(@dkim_selector, dkim_selector),
    dkim_public_key = COALESCE(
        @dkim_public_key,
        dkim_public_key
    ),
    dkim_private_key = COALESCE(
        @dkim_private_key,
        dkim_private_key
    ),
    spf_record = COALESCE(@spf_record, spf_record),
    verified = COALESCE(@verified, verified),
    mx_record = COALESCE(@mx_record, mx_record),
    updated_at = CURRENT_TIMESTAMP
WHERE
    id = @id
    AND company_id = @company_id;

-- name: GetDomainByIDAndCompany :one
SELECT *
FROM domains
WHERE
    id = @id
    AND company_id = @company_id
    AND deleted_at IS NULL;

-- name: ListDomainsByCompany :many
SELECT *
FROM domains
WHERE
    company_id = @company_id
    AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT @rowlimit
OFFSET
    @rowoffset;

-- name: SoftDeleteDomain :exec
UPDATE domains
SET
    deleted_at = now(),
    updated_at = now()
WHERE
    company_id = $1
    AND id = $2;

-- name: CountDomainByCompanyID :one
SELECT COUNT(*)
FROM domains
WHERE
    company_id = $1
    AND deleted_at IS NULL;

-- name: FindDomainByNameAndCompany :one
SELECT *
FROM domains
WHERE
    domain = $1
    AND company_id = $2
    AND deleted_at IS NULL;

-- name: SoftDeleteSender :exec
UPDATE senders
SET
    deleted_at = now(),
    updated_at = now()
WHERE
    id = $1
    AND company_id = $2;

-- name: ListSendersByCompanyId :many
SELECT *
FROM senders
WHERE
    company_id = $1
    AND deleted_at IS NULL
ORDER BY created_at DESC
LIMIT $2
OFFSET
    $3;

-- name: CountSenderByCompanyID :one
SELECT COUNT(*)
FROM senders
WHERE
    company_id = $1
    AND deleted_at IS NULL;

-- name: UpdateSender :one
UPDATE senders
SET
    user_id = COALESCE($3, user_id),
    name = COALESCE($4, name),
    email = COALESCE($5, email),
    verified = COALESCE($6, verified),
    is_signed = COALESCE($7, is_signed),
    domain_id = COALESCE($8, domain_id),
    updated_at = NOW()
WHERE
    id = $1
    AND company_id = $2 RETURNING *;

-- name: UpdateSenderVerified :exec
UPDATE senders
SET
    verified = $1,
    updated_at = NOW()
WHERE
    company_id = $2
    AND email = $3;

-- name: GetSenderById :one
SELECT *
FROM senders
WHERE
    company_id = $1
    AND id = $2
    AND deleted_at IS NULL;
