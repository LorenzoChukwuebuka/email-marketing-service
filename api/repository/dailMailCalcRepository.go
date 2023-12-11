package repository

import (
	"database/sql"
	"email-marketing-service/api/model"
	"fmt"
	"time"
)

type DailyMailCalcRepository struct {
	DB *sql.DB
}

func NewDailyMailCalcRepository(db *sql.DB) *DailyMailCalcRepository {
	return &DailyMailCalcRepository{DB: db}
}

func (r *DailyMailCalcRepository) CreateRecordDailyMailCalculation(d *model.DailyMailCalcModel) error {
	query := `
		INSERT INTO daily_mail_calc (
			uuid, subscription_id, mails_for_a_day, mails_sent, created_at,mails_remaining
		) VALUES ($1, $2, $3, $4, $5,$6)
	`

	stmt, err := r.DB.Prepare(query)
	if err != nil {
		fmt.Println("Error preparing statement:", err)
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(
		d.UUID,
		d.SubscriptionID,
		d.MailsForADay,
		d.MailsSent,
		d.CreatedAt,
		d.RemainingMails,
	)

	if err != nil {
		fmt.Println("Error executing statement:", err)
		return err
	}

	return nil
}

func (r *DailyMailCalcRepository) GetDailyMailRecordForToday(userId int) (*model.DailyMailCalcResponseModel, error) {
	// Get the current date in the format "YYYY-MM-DD"
	currentDate := time.Now().Format("2006-01-02")

	// Prepare the SQL statement
	query := `
		SELECT id, uuid, subscription_id, mails_for_a_day, mails_sent, created_at, updated_at,mails_remaining
		FROM public.daily_mail_calc
		WHERE subscription_id = $1 AND DATE(created_at) = $2
		LIMIT 1
	`

	// Execute the query
	row := r.DB.QueryRow(query, userId, currentDate)

	var updatedAt sql.NullTime
	// Parse the result
	var dailyMailRecord model.DailyMailCalcResponseModel
	err := row.Scan(
		&dailyMailRecord.ID,
		&dailyMailRecord.UUID,
		&dailyMailRecord.SubscriptionID,
		&dailyMailRecord.MailsForADay,
		&dailyMailRecord.MailsSent,
		&dailyMailRecord.CreatedAt,
		&updatedAt,
		&dailyMailRecord.RemainingMails,
	)

	SetTime(updatedAt, &dailyMailRecord.UpdatedAt)

	if err != nil {
		if err == sql.ErrNoRows {
			// No records found
			return nil, nil
		}
		fmt.Println("Error scanning row:", err)
		return nil, err
	}

	return &dailyMailRecord, nil
}

func (r *DailyMailCalcRepository) UpdateDailyMailCalcRepository(d *model.DailyMailCalcModel) error {

	query := "Update daily_mail_calc SET mails_sent = $1,mails_remaining = $2 Where uuid = $3"

	stmt, err := r.DB.Prepare(query)

	if err != nil {
		return err
	}

	defer stmt.Close()

	_, err = stmt.Exec(d.MailsSent, d.RemainingMails, d.UUID)

	if err != nil {
		return err
	}
	return nil
}
