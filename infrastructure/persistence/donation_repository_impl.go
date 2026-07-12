package persistence

import (
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

// DonationModel is the GORM model for database
type DonationModel struct {
	ID            uint       `gorm:"primaryKey"`
	Type          string     `gorm:"size:50;not null;index"` // donate or sponsor
	DonorType     string     `gorm:"size:50;not null;index"` // user or company
	UserID        *uint      `gorm:"index"`
	SupporterID   *uint      `gorm:"index"`
	CompanyName   string     `gorm:"size:255"`
	CompanyEmail  string     `gorm:"size:255"`
	CompanyPhone  string     `gorm:"size:50"`
	Amount        float64    `gorm:"not null"`
	Currency      string     `gorm:"size:10;not null;default:'USD'"`
	Message       string     `gorm:"type:text"`
	Status        string     `gorm:"size:50;not null;default:'pending';index"`
	TransactionID string     `gorm:"size:255;index"`
	PaymentMethod string     `gorm:"size:100"`
	QRExpiresAt   *time.Time `gorm:"index"`
	EventID       *uint      `gorm:"index"`
	CreatedAt     time.Time  `gorm:"autoCreateTime"`
	UpdatedAt     time.Time  `gorm:"autoUpdateTime"`

	// Relations
	User      UserModel      `gorm:"foreignKey:UserID"`
	Supporter SupporterModel `gorm:"foreignKey:SupporterID"`
	Event     EventModel     `gorm:"foreignKey:EventID"`
}

func (DonationModel) TableName() string {
	return "donations"
}

type donationRepositoryImpl struct {
	db *gorm.DB
}

func NewDonationRepository(db *gorm.DB) repository.DonationRepository {
	return &donationRepositoryImpl{db: db}
}

func (r *donationRepositoryImpl) Save(donation *entity.Donation) error {
	model := r.entityToModel(donation)
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}

	donation.ID = model.ID
	donation.CreatedAt = model.CreatedAt
	donation.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *donationRepositoryImpl) FindByID(id uint) (*entity.Donation, error) {
	var model DonationModel
	if err := r.db.Preload("User").Preload("Supporter").Preload("Event").First(&model, id).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *donationRepositoryImpl) FindAll(limit, offset int) ([]*entity.Donation, error) {
	var models []DonationModel
	query := r.db.Preload("User").Preload("Supporter").Preload("Event").Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *donationRepositoryImpl) FindByUserID(userID uint, limit, offset int) ([]*entity.Donation, error) {
	var models []DonationModel
	query := r.db.Preload("User").Preload("Supporter").Preload("Event").
		Where("user_id = ?", userID).
		Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *donationRepositoryImpl) FindByType(donationType entity.DonationType, limit, offset int) ([]*entity.Donation, error) {
	var models []DonationModel
	query := r.db.Preload("User").Preload("Supporter").Preload("Event").
		Where("type = ?", string(donationType)).
		Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *donationRepositoryImpl) FindByDonorType(donorType entity.DonorType, limit, offset int) ([]*entity.Donation, error) {
	var models []DonationModel
	query := r.db.Preload("User").Preload("Supporter").Preload("Event").
		Where("donor_type = ?", string(donorType)).
		Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *donationRepositoryImpl) FindByEventID(eventID uint, limit, offset int) ([]*entity.Donation, error) {
	var models []DonationModel
	query := r.db.Preload("User").Preload("Supporter").Preload("Event").
		Where("event_id = ?", eventID).
		Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *donationRepositoryImpl) FindByStatus(status entity.DonationStatus, limit, offset int) ([]*entity.Donation, error) {
	var models []DonationModel
	query := r.db.Preload("User").Preload("Supporter").Preload("Event").
		Where("status = ?", string(status)).
		Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit)
	}
	if offset > 0 {
		query = query.Offset(offset)
	}
	
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *donationRepositoryImpl) Update(donation *entity.Donation) error {
	model := r.entityToModel(donation)
	return r.db.Save(&model).Error
}

func (r *donationRepositoryImpl) UpdateStatus(id uint, status entity.DonationStatus, transactionID, paymentMethod string) error {
	return r.db.Model(&DonationModel{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":         string(status),
		"transaction_id": transactionID,
		"payment_method": paymentMethod,
		"updated_at":     time.Now(),
	}).Error
}

func (r *donationRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&DonationModel{}, id).Error
}

func (r *donationRepositoryImpl) GetTotalAmount() (float64, error) {
	var total float64
	err := r.db.Model(&DonationModel{}).
		Where("status = ?", string(entity.DonationStatusCompleted)).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *donationRepositoryImpl) GetTotalAmountByType(donationType entity.DonationType) (float64, error) {
	var total float64
	err := r.db.Model(&DonationModel{}).
		Where("type = ? AND status = ?", string(donationType), string(entity.DonationStatusCompleted)).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *donationRepositoryImpl) GetTotalAmountByEventID(eventID uint) (float64, error) {
	var total float64
	err := r.db.Model(&DonationModel{}).
		Where("event_id = ? AND status = ?", eventID, string(entity.DonationStatusCompleted)).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	return total, err
}

func (r *donationRepositoryImpl) Count() (int64, error) {
	var count int64
	err := r.db.Model(&DonationModel{}).Count(&count).Error
	return count, err
}

func (r *donationRepositoryImpl) CountByType(donationType entity.DonationType) (int64, error) {
	var count int64
	err := r.db.Model(&DonationModel{}).
		Where("type = ?", string(donationType)).
		Count(&count).Error
	return count, err
}

func (r *donationRepositoryImpl) entityToModel(donation *entity.Donation) *DonationModel {
	return &DonationModel{
		ID:            donation.ID,
		Type:          string(donation.Type),
		DonorType:     string(donation.DonorType),
		UserID:        donation.UserID,
		SupporterID:   donation.SupporterID,
		CompanyName:   donation.CompanyName,
		CompanyEmail:  donation.CompanyEmail,
		CompanyPhone:  donation.CompanyPhone,
		Amount:        donation.Amount,
		Currency:      donation.Currency,
		Message:       donation.Message,
		Status:        string(donation.Status),
		TransactionID: donation.TransactionID,
		PaymentMethod: donation.PaymentMethod,
		QRExpiresAt:   donation.QRExpiresAt,
		EventID:       donation.EventID,
		CreatedAt:     donation.CreatedAt,
		UpdatedAt:     donation.UpdatedAt,
	}
}

func (r *donationRepositoryImpl) modelToEntity(model *DonationModel) *entity.Donation {
	donation := &entity.Donation{
		ID:            model.ID,
		Type:          entity.DonationType(model.Type),
		DonorType:     entity.DonorType(model.DonorType),
		UserID:        model.UserID,
		SupporterID:   model.SupporterID,
		CompanyName:   model.CompanyName,
		CompanyEmail:  model.CompanyEmail,
		CompanyPhone:  model.CompanyPhone,
		Amount:        model.Amount,
		Currency:      model.Currency,
		Message:       model.Message,
		Status:        entity.DonationStatus(model.Status),
		TransactionID: model.TransactionID,
		PaymentMethod: model.PaymentMethod,
		QRExpiresAt:   model.QRExpiresAt,
		EventID:       model.EventID,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
	}

	// Convert user if loaded
	if model.User.ID != 0 {
		donation.User = &entity.User{
			ID:        model.User.ID,
			Username:  model.User.Username,
			Name:  model.User.Name,
			Profile:   model.User.Profile,
			Email:     model.User.Email,
			Role:      model.User.Role,
			CreatedAt: model.User.CreatedAt,
			UpdatedAt: model.User.UpdatedAt,
		}
	}

	// Convert supporter if loaded
	if model.Supporter.ID != 0 {
		donation.Supporter = &entity.Supporter{
			ID:          model.Supporter.ID,
			Name:        model.Supporter.Name,
			Email:       model.Supporter.Email,
			Phone:       model.Supporter.Phone,
			Type:        entity.SupporterType(model.Supporter.Type),
			Website:     model.Supporter.Website,
			Address:     model.Supporter.Address,
			Logo:        model.Supporter.Logo,
			Description: model.Supporter.Description,
			UserID:      model.Supporter.UserID,
			CreatedAt:   model.Supporter.CreatedAt,
			UpdatedAt:   model.Supporter.UpdatedAt,
		}
	}

	// Convert event if loaded
	if model.Event.ID != 0 {
		donation.Event = &entity.Event{
			ID:        model.Event.ID,
			Title:     model.Event.Title,
			Content:   model.Event.Content,
			Cover:     model.Event.Cover,
			Location:  model.Event.Location,
			StartTime: model.Event.StartTime,
			EndTime:   model.Event.EndTime,
			UserID:    model.Event.UserID,
			CreatedAt: model.Event.CreatedAt,
			UpdatedAt: model.Event.UpdatedAt,
		}
	}

	return donation
}

func (r *donationRepositoryImpl) modelsToEntities(models []DonationModel) []*entity.Donation {
	entities := make([]*entity.Donation, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(&model)
	}
	return entities
}

