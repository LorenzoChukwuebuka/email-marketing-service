package repository

import (
	"email-marketing-service/api/model"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type DailyMailCalcRepository struct {
	DB *gorm.DB
}

func NewDailyMailCalcRepository(db *gorm.DB) *DailyMailCalcRepository {
	return &DailyMailCalcRepository{DB: db}
}

func (r *DailyMailCalcRepository) CreateRecordDailyMailCalculation(d *model.DailyMailCalc) error {
	if err := r.DB.Create(&d).Error; err != nil {
		return fmt.Errorf("failed to insert plan: %w", err)
	}
	return nil
}

func (r *DailyMailCalcRepository) GetDailyMailRecordForToday(userId int) (*model.DailyMailCalcResponseModel, error) {
	var record model.DailyMailCalc

	today := time.Now().Format("2006-01-02")
	if err := r.DB.Where("user_id = ? AND  created_at >= ? AND created_at <  ?", userId, today, today+" 23:59:59").First(&record).Error; err != nil {
		return nil, err
	}

	println(&record)

	return nil, nil
	//return &record, nil

}

func (r *DailyMailCalcRepository) UpdateDailyMailCalcRepository(d *model.DailyMailCalc) error {

	return nil
}
