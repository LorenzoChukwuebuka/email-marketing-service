-- name: CreateDailyEmailUsage :one
INSERT INTO
    email_usage (
        company_id,
        subscription_id,
        usage_period_start,
        usage_period_end,
        period_type,
        emails_sent,
        emails_limit
    )
VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING *;

-- name: GetCurrentBillingPeriod :one
SELECT 
    s.company_id,
    s.id as subscription_id,
    s.created_at as subscription_start,
    DATE_TRUNC('month', CURRENT_DATE)::DATE as current_period_start,
    (DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month - 1 day')::DATE as current_period_end,
    ml.daily_limit,
    ml.monthly_limit,
    EXTRACT(DAY FROM (DATE_TRUNC('month', CURRENT_DATE) + INTERVAL '1 month - 1 day') - DATE_TRUNC('month', CURRENT_DATE)) + 1 as days_in_period
FROM subscriptions s
JOIN plans p ON s.plan_id = p.id
JOIN mailing_limits ml ON p.id = ml.plan_id
WHERE s.company_id = $1 
  AND s.status = 'active';


-- name: GetEmailUsageByCompanyAndPeriod :one
SELECT *
FROM email_usage
WHERE
    company_id = $1
    AND usage_period_start = $2
    AND period_type = $3;

-- name: GetCurrentEmailUsage :one
SELECT *
FROM email_usage
WHERE
    company_id = $1
    AND usage_period_start <= CURRENT_DATE
    AND usage_period_end >= CURRENT_DATE;

-- name: GetEmailUsageByCompany :many
SELECT *
FROM email_usage
WHERE
    company_id = $1
ORDER BY usage_period_start DESC;

-- name: GetEmailUsageBySubscription :many
SELECT *
FROM email_usage
WHERE
    subscription_id = $1
ORDER BY usage_period_start DESC;

-- name: GetEmailUsageInDateRange :many
SELECT *
FROM email_usage
WHERE
    company_id = $1
    AND usage_period_start >= $2
    AND usage_period_end <= $3
    AND period_type = $4
ORDER BY usage_period_start ASC;

-- name: IncrementEmailsSent :one
UPDATE email_usage
SET
    emails_sent = emails_sent + $3,
    updated_at = CURRENT_TIMESTAMP
WHERE
    company_id = $1
    AND usage_period_start = $2
    AND period_type = 'daily' RETURNING *;

-- name: GetEmailUsageStats :one
SELECT
    company_id,
    period_type,
    COUNT(*) as total_periods,
    SUM(emails_sent) as total_emails_sent,
    AVG(emails_sent) as avg_emails_per_period,
    MAX(emails_sent) as max_emails_in_period,
    SUM(emails_limit) as total_email_limits
FROM email_usage
WHERE
    company_id = $1
    AND period_type = $2
GROUP BY
    company_id,
    period_type;

-- name: UpdateEmailsSentAndRemaining :one
UPDATE email_usage
SET
    emails_sent = emails_sent + $2,
    remaining_emails = remaining_emails - $2,
    updated_at = CURRENT_TIMESTAMP
WHERE
    company_id = $1
    AND id = $3 RETURNING *;

-- name: GetCompaniesNearLimit :many
SELECT 
    eu.*,
    c.companyname as company_name,
    ROUND((eu.emails_sent::DECIMAL / eu.emails_limit::DECIMAL) * 100, 2) as usage_percentage
FROM email_usage eu
JOIN companies c ON c.id = eu.company_id
WHERE eu.period_type = $1
  AND eu.usage_period_start <= CURRENT_DATE 
  AND eu.usage_period_end >= CURRENT_DATE
  AND (eu.emails_sent::DECIMAL / eu.emails_limit::DECIMAL) >= $2
ORDER BY usage_percentage DESC;

-- name: GetMonthlyEmailTrends :many
SELECT
    DATE_TRUNC ('month', usage_period_start) as month,
    SUM(emails_sent) as total_emails_sent,
    AVG(emails_sent) as avg_daily_emails,
    COUNT(*) as active_days
FROM email_usage
WHERE
    company_id = $1
    AND period_type = 'daily'
    AND usage_period_start >= $2
    AND usage_period_start <= $3
GROUP BY
    DATE_TRUNC ('month', usage_period_start)
ORDER BY month ASC;

-- name: CheckEmailLimitExceeded :one
SELECT
    *,
    CASE
        WHEN emails_sent >= emails_limit THEN true
        ELSE false
    END as limit_exceeded,
    (emails_limit - emails_sent) as remaining_emails
FROM email_usage
WHERE
    company_id = $1
    AND usage_period_start = $2
    AND period_type = $3;