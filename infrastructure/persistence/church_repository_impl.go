package persistence

import (
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

type ChurchModel struct {
	ID              uint       `gorm:"primaryKey"`
	Name        string     `gorm:"size:255;not null;uniqueIndex"`
	Address         string     `gorm:"type:text"`
	Phone           string     `gorm:"size:50"`
	Email           string     `gorm:"size:255"`
	Website         string     `gorm:"size:255"`
	PastorName      string     `gorm:"size:255"`
	Description     string     `gorm:"type:text"`
	Logo            string     `gorm:"size:255"`
	EstablishedDate *time.Time
	Denomination    string    `gorm:"size:100"`
	OwnerID         *uint     `gorm:"index"`
	CreatedAt       time.Time `gorm:"autoCreateTime"`
	UpdatedAt       time.Time `gorm:"autoUpdateTime"`
}

func (ChurchModel) TableName() string {
	return "churches"
}

type churchRepositoryImpl struct {
	db *gorm.DB
}

func NewChurchRepository(db *gorm.DB) repository.ChurchRepository {
	return &churchRepositoryImpl{db: db}
}

func (r *churchRepositoryImpl) Create(church *entity.Church) error {
	model := r.entityToModel(church)
	if err := r.db.Create(model).Error; err != nil {
		return err
	}
	*church = *r.modelToEntity(model)
	return nil
}

func (r *churchRepositoryImpl) FindByID(id uint) (*entity.Church, error) {
	var model ChurchModel
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	
	// Load owner separately if needed
	if model.OwnerID != nil {
		var ownerModel UserModel
		if err := r.db.First(&ownerModel, model.OwnerID).Error; err == nil {
			church := r.modelToEntity(&model)
			church.Owner = &entity.User{
				ID:       ownerModel.ID,
				Username: ownerModel.Username,
				Name: ownerModel.Name,
				Email:    ownerModel.Email,
				Profile:  ownerModel.Profile,
			}
			return church, nil
		}
	}
	
	return r.modelToEntity(&model), nil
}

func (r *churchRepositoryImpl) FindByName(name string) (*entity.Church, error) {
	var model ChurchModel
	if err := r.db.Where("name = ?", name).First(&model).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *churchRepositoryImpl) FindAll(limit, offset int) ([]*entity.Church, error) {
	var models []ChurchModel
	query := r.db.Order("name ASC")
	
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	
	return r.modelsToEntities(models), nil
}

func (r *churchRepositoryImpl) FindByDenomination(denomination string, limit, offset int) ([]*entity.Church, error) {
	var models []ChurchModel
	query := r.db.Where("denomination = ?", denomination).Order("name ASC")
	
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	
	return r.modelsToEntities(models), nil
}

func (r *churchRepositoryImpl) Update(church *entity.Church) error {
	model := r.entityToModel(church)
	return r.db.Save(model).Error
}

func (r *churchRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&ChurchModel{}, id).Error
}

func (r *churchRepositoryImpl) Count() (int64, error) {
	var count int64
	err := r.db.Model(&ChurchModel{}).Count(&count).Error
	return count, err
}

func (r *churchRepositoryImpl) CountMembers(churchID uint, status string) (int64, error) {
	var count int64
	query := r.db.Model(&UserModel{}).Where("church_id = ?", churchID)
	if status != "" {
		query = query.Where("church_status = ?", status)
	}
	err := query.Count(&count).Error
	return count, err
}

func (r *churchRepositoryImpl) FindMembers(churchID uint, status string, limit, offset int) ([]*entity.User, error) {
	var models []UserModel
	query := r.db.Where("church_id = ?", churchID)
	if status != "" {
		query = query.Where("church_status = ?", status)
	}
	
	if limit > 0 {
		query = query.Limit(limit).Offset(offset)
	}
	
	if err := query.Find(&models).Error; err != nil {
		return nil, err
	}
	
	users := make([]*entity.User, len(models))
	for i, model := range models {
		users[i] = &entity.User{
			ID:           model.ID,
			Username:     model.Username,
			Name:     model.Name,
			Email:        model.Email,
			Profile:      model.Profile,
			ChurchID:     model.ChurchID,
			ChurchStatus: entity.ChurchMembershipStatus(model.ChurchStatus),
			Birthday:     model.Birthday,
			Bio:          model.Bio,
			CreatedAt:    model.CreatedAt,
			UpdatedAt:    model.UpdatedAt,
		}
	}
	return users, nil
}

// Helper methods for conversion
func (r *churchRepositoryImpl) entityToModel(church *entity.Church) *ChurchModel {
	return &ChurchModel{
		ID:              church.ID,
		Name:        church.Name,
		Address:         church.Address,
		Phone:           church.Phone,
		Email:           church.Email,
		Website:         church.Website,
		PastorName:      church.PastorName,
		Description:     church.Description,
		Logo:            church.Logo,
		EstablishedDate: church.EstablishedDate,
		Denomination:    church.Denomination,
		OwnerID:         church.OwnerID,
		CreatedAt:       church.CreatedAt,
		UpdatedAt:       church.UpdatedAt,
	}
}

func (r *churchRepositoryImpl) modelToEntity(model *ChurchModel) *entity.Church {
	return &entity.Church{
		ID:              model.ID,
		Name:        	 model.Name,
		Address:         model.Address,
		Phone:           model.Phone,
		Email:           model.Email,
		Website:         model.Website,
		PastorName:      model.PastorName,
		Description:     model.Description,
		Logo:            model.Logo,
		EstablishedDate: model.EstablishedDate,
		Denomination:    model.Denomination,
		OwnerID:         model.OwnerID,
		CreatedAt:       model.CreatedAt,
		UpdatedAt:       model.UpdatedAt,
	}
}

func (r *churchRepositoryImpl) modelsToEntities(models []ChurchModel) []*entity.Church {
	churches := make([]*entity.Church, len(models))
	for i, model := range models {
		churches[i] = r.modelToEntity(&model)
	}
	return churches
}

