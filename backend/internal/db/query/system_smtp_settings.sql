-- name: CreateSystemsSMTPSettings :one
INSERT INTO
    systems_smtp_settings (
        txt_record,
        dmarc_record,
        dkim_selector,
        dkim_public_key,
        dkim_private_key,
        spf_record,
        verified,
        mx_record,
        domain
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
        $9
    ) RETURNING *;


-- name: DeleteSystemsSMTPSetting :exec
DELETE FROM systems_smtp_settings WHERE domain = $1;

-- name: GetSMTPSettingByDomain :one
SELECT * FROM systems_smtp_settings
WHERE domain = $1;