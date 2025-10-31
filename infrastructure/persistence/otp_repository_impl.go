package persistence

import (
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

// OTPModel is the GORM model for database
type OTPModel struct {
	ID        uint      `gorm:"primaryKey;autoIncrement"`
	Email     string    `gorm:"size:100;not null;index"`
	Code      string    `gorm:"size:10;not null"`
	Purpose   string    `gorm:"size:50;not null;index"`
	ExpiresAt time.Time `gorm:"not null;index"`
	Verified  bool      `gorm:"default:false"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (OTPModel) TableName() string {
	return "otps"
}

type otpRepositoryImpl struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) repository.OTPRepository {
	return &otpRepositoryImpl{db: db}
}

func (r *otpRepositoryImpl) Save(otp *entity.OTP) error {
	model := r.entityToModel(otp)
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	otp.ID = model.ID
	otp.CreatedAt = model.CreatedAt
	otp.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *otpRepositoryImpl) FindByEmailAndPurpose(email string, purpose entity.OTPPurpose) (*entity.OTP, error) {
	var model OTPModel
	if err := r.db.Where("email = ? AND purpose = ? AND verified = false AND expires_at > ?", 
		email, purpose, time.Now()).
		Order("created_at DESC").
		First(&model).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *otpRepositoryImpl) FindByEmailCodeAndPurpose(email, code string, purpose entity.OTPPurpose) (*entity.OTP, error) {
	var model OTPModel
	if err := r.db.Where("email = ? AND code = ? AND purpose = ?", email, code, purpose).
		Order("created_at DESC").
		First(&model).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *otpRepositoryImpl) Update(otp *entity.OTP) error {
	model := r.entityToModel(otp)
	return r.db.Save(&model).Error
}

func (r *otpRepositoryImpl) DeleteByEmail(email string) error {
	return r.db.Where("email = ?", email).Delete(&OTPModel{}).Error
}

func (r *otpRepositoryImpl) DeleteByEmailAndPurpose(email string, purpose entity.OTPPurpose) error {
	return r.db.Where("email = ? AND purpose = ?", email, purpose).Delete(&OTPModel{}).Error
}

func (r *otpRepositoryImpl) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&OTPModel{}).Error
}

func (r *otpRepositoryImpl) entityToModel(otp *entity.OTP) *OTPModel {
	return &OTPModel{
		ID:        otp.ID,
		Email:     otp.Email,
		Code:      otp.Code,
		Purpose:   string(otp.Purpose),
		ExpiresAt: otp.ExpiresAt,
		Verified:  otp.Verified,
		CreatedAt: otp.CreatedAt,
		UpdatedAt: otp.UpdatedAt,
	}
}

func (r *otpRepositoryImpl) modelToEntity(model *OTPModel) *entity.OTP {
	return &entity.OTP{
		ID:        model.ID,
		Email:     model.Email,
		Code:      model.Code,
		Purpose:   entity.OTPPurpose(model.Purpose),
		ExpiresAt: model.ExpiresAt,
		Verified:  model.Verified,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

