package repository

import (
	"database/sql"
	"email-marketing-service/api/model"
	"fmt"
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
			uuid, subscription_id, mails_for_a_day, mails_sent, created_at, updated_at
		) VALUES ($1, $2, $3, $4, $5, $6)
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
		d.UpdatedAt,
	)

	if err != nil {
		fmt.Println("Error executing statement:", err)
		return err
	}

	return nil
}
