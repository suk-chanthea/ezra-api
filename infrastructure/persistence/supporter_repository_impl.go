package persistence

import (
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

type SupporterModel struct {
	ID          uint      `gorm:"primaryKey"`
	Name        string    `gorm:"size:255;not null"`
	Email       string    `gorm:"size:255;not null;uniqueIndex"`
	Phone       string    `gorm:"size:50"`
	Type        string    `gorm:"size:50;not null;default:'company'"`
	Website     string    `gorm:"size:255"`
	Address     string    `gorm:"type:text"`
	Logo        string    `gorm:"size:255"`
	Description string    `gorm:"type:text"`
	UserID      *uint     `gorm:"index"`
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	// Relations
	User      UserModel      `gorm:"foreignKey:UserID"`
	Donations []DonationModel `gorm:"foreignKey:SupporterID"`
}

func (SupporterModel) TableName() string {
	return "supporters"
}

type supporterRepositoryImpl struct {
	db *gorm.DB
}

func NewSupporterRepository(db *gorm.DB) repository.SupporterRepository {
	return &supporterRepositoryImpl{db: db}
}

func (r *supporterRepositoryImpl) Create(supporter *entity.Supporter) error {
	model := r.entityToModel(supporter)
	if err := r.db.Create(model).Error; err != nil {
		return err
	}
	*supporter = *r.modelToEntity(model)
	return nil
}

func (r *supporterRepositoryImpl) FindByID(id uint) (*entity.Supporter, error) {
	var model SupporterModel
	if err := r.db.Preload("User").Preload("Donations").First(&model, id).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *supporterRepositoryImpl) FindByEmail(email string) (*entity.Supporter, error) {
	var model SupporterModel
	if err := r.db.Where("email = ?", email).Preload("User").First(&model).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *supporterRepositoryImpl) FindAll(limit, offset int) ([]*entity.Supporter, error) {
	var models []SupporterModel
	query := r.db.Preload("User").Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	
	return r.modelsToEntities(models), nil
}

func (r *supporterRepositoryImpl) FindByType(supporterType entity.SupporterType, limit, offset int) ([]*entity.Supporter, error) {
	var models []SupporterModel
	query := r.db.Where("type = ?", string(supporterType)).Preload("User").Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	
	return r.modelsToEntities(models), nil
}

func (r *supporterRepositoryImpl) FindByUser(userID uint, limit, offset int) ([]*entity.Supporter, error) {
	var models []SupporterModel
	query := r.db.Where("user_id = ?", userID).Preload("User").Order("created_at DESC")
	
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	
	return r.modelsToEntities(models), nil
}

func (r *supporterRepositoryImpl) Update(supporter *entity.Supporter) error {
	model := r.entityToModel(supporter)
	return r.db.Save(model).Error
}

func (r *supporterRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&SupporterModel{}, id).Error
}

func (r *supporterRepositoryImpl) Count() (int64, error) {
	var count int64
	err := r.db.Model(&SupporterModel{}).Count(&count).Error
	return count, err
}

func (r *supporterRepositoryImpl) CountByType(supporterType entity.SupporterType) (int64, error) {
	var count int64
	err := r.db.Model(&SupporterModel{}).Where("type = ?", string(supporterType)).Count(&count).Error
	return count, err
}

// Helper methods for conversion
func (r *supporterRepositoryImpl) entityToModel(supporter *entity.Supporter) *SupporterModel {
	return &SupporterModel{
		ID:          supporter.ID,
		Name:        supporter.Name,
		Email:       supporter.Email,
		Phone:       supporter.Phone,
		Type:        string(supporter.Type),
		Website:     supporter.Website,
		Address:     supporter.Address,
		Logo:        supporter.Logo,
		Description: supporter.Description,
		UserID:      supporter.UserID,
		CreatedAt:   supporter.CreatedAt,
		UpdatedAt:   supporter.UpdatedAt,
	}
}

func (r *supporterRepositoryImpl) modelToEntity(model *SupporterModel) *entity.Supporter {
	supporter := &entity.Supporter{
		ID:          model.ID,
		Name:        model.Name,
		Email:       model.Email,
		Phone:       model.Phone,
		Type:        entity.SupporterType(model.Type),
		Website:     model.Website,
		Address:     model.Address,
		Logo:        model.Logo,
		Description: model.Description,
		UserID:      model.UserID,
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
	}

	// Convert related user if loaded
	if model.User.ID != 0 {
		supporter.User = &entity.User{
			ID:       model.User.ID,
			Username: model.User.Username,
			Name: model.User.Name,
			Email:    model.User.Email,
			Profile:  model.User.Profile,
			Role:     model.User.Role,
		}
	}

	return supporter
}

func (r *supporterRepositoryImpl) modelsToEntities(models []SupporterModel) []*entity.Supporter {
	supporters := make([]*entity.Supporter, len(models))
	for i, model := range models {
		supporters[i] = r.modelToEntity(&model)
	}
	return supporters
}

