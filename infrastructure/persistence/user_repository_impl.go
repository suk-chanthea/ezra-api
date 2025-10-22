package persistence

import (
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

// UserModel is the GORM model for database
type UserModel struct {
	ID         uint      `gorm:"primaryKey;autoIncrement"`
	Username   string    `gorm:"size:100;not null"`
	Fullname   string    `gorm:"size:100;not null"`
	Profile    string    `gorm:"size:255"`
	Email      string    `gorm:"size:100;unique;not null"`
	Password   string    `gorm:"size:255"`
	Role       string    `gorm:"size:20;default:user"`
	Token      string    `gorm:"size:255"`
	Provider   string    `gorm:"size:50;default:local"`
	ProviderID string    `gorm:"size:255;index"`
	CreatedAt  time.Time `gorm:"autoCreateTime"`
	UpdatedAt  time.Time `gorm:"autoUpdateTime"`
}

func (UserModel) TableName() string {
	return "users"
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) repository.UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) Save(user *entity.User) error {
	model := r.entityToModel(user)
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	user.ID = model.ID
	user.CreatedAt = model.CreatedAt
	user.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *userRepositoryImpl) FindByUsername(username string) (*entity.User, error) {
	var model UserModel
	if err := r.db.Where("username = ?", username).First(&model).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *userRepositoryImpl) FindByEmail(email string) (*entity.User, error) {
	var model UserModel
	if err := r.db.Where("email = ?", email).First(&model).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *userRepositoryImpl) FindByID(id uint) (*entity.User, error) {
	var model UserModel
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *userRepositoryImpl) Update(user *entity.User) error {
	model := r.entityToModel(user)
	return r.db.Save(&model).Error
}

func (r *userRepositoryImpl) FindByProviderID(provider, providerID string) (*entity.User, error) {
	var model UserModel
	if err := r.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&model).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *userRepositoryImpl) UpdateToken(id uint, token string) error {
	return r.db.Model(&UserModel{}).
		Where("id = ?", id).
		Update("token", token).Error
}

func (r *userRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&UserModel{}, id).Error
}

func (r *userRepositoryImpl) entityToModel(user *entity.User) *UserModel {
	return &UserModel{
		ID:         user.ID,
		Username:   user.Username,
		Fullname:   user.Fullname,
		Profile:    user.Profile,
		Email:      user.Email,
		Password:   user.Password,
		Role:       user.Role,
		Token:      user.Token,
		Provider:   user.Provider,
		ProviderID: user.ProviderID,
		CreatedAt:  user.CreatedAt,
		UpdatedAt:  user.UpdatedAt,
	}
}

func (r *userRepositoryImpl) modelToEntity(model *UserModel) *entity.User {
	return &entity.User{
		ID:         model.ID,
		Username:   model.Username,
		Fullname:   model.Fullname,
		Profile:    model.Profile,
		Email:      model.Email,
		Password:   model.Password,
		Role:       model.Role,
		Token:      model.Token,
		Provider:   model.Provider,
		ProviderID: model.ProviderID,
		CreatedAt:  model.CreatedAt,
		UpdatedAt:  model.UpdatedAt,
	}
}