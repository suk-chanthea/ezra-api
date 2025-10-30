package repository

import "github.com/suk-chanthea/ezra/domain/entity"

// DonationRepository defines the interface for donation data operations
type DonationRepository interface {
	Save(donation *entity.Donation) error
	FindByID(id uint) (*entity.Donation, error)
	FindAll(limit, offset int) ([]*entity.Donation, error)
	FindByUserID(userID uint, limit, offset int) ([]*entity.Donation, error)
	FindByType(donationType entity.DonationType, limit, offset int) ([]*entity.Donation, error)
	FindByDonorType(donorType entity.DonorType, limit, offset int) ([]*entity.Donation, error)
	FindByEventID(eventID uint, limit, offset int) ([]*entity.Donation, error)
	FindByStatus(status entity.DonationStatus, limit, offset int) ([]*entity.Donation, error)
	Update(donation *entity.Donation) error
	UpdateStatus(id uint, status entity.DonationStatus, transactionID, paymentMethod string) error
	Delete(id uint) error
	GetTotalAmount() (float64, error)
	GetTotalAmountByType(donationType entity.DonationType) (float64, error)
	GetTotalAmountByEventID(eventID uint) (float64, error)
	Count() (int64, error)
	CountByType(donationType entity.DonationType) (int64, error)
}

