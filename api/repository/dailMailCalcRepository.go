package repository

import (
	"email-marketing-service/api/model"
	"gorm.io/gorm"
)

type DailyMailCalcRepository struct {
	DB *gorm.DB
}

func NewDailyMailCalcRepository(db *gorm.DB) *DailyMailCalcRepository {
	return &DailyMailCalcRepository{DB: db}
}

func (r *DailyMailCalcRepository) CreateRecordDailyMailCalculation(d *model.DailyMailCalc) error {

	return nil
}

func (r *DailyMailCalcRepository) GetDailyMailRecordForToday(userId int) (*model.DailyMailCalcResponseModel, error) {

	return nil, nil
}

func (r *DailyMailCalcRepository) UpdateDailyMailCalcRepository(d *model.DailyMailCalc) error {

	return nil
}
