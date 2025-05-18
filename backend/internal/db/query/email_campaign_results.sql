-- name: MarkEmailAsDelivered :exec
UPDATE email_campaign_results
SET
    sent_at = NOW(),
    bounce_status = '',
    updated_at = NOW()
WHERE
    recipient_email = $1
    AND deleted_at IS NULL;

-- name: UpdateBounceStatus :exec
UPDATE email_campaign_results
SET
    bounce_status = $1,
    updated_at = NOW()
WHERE
    recipient_email = $2
    AND deleted_at IS NULL;