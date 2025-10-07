package persistence

import (
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

// EventModel is the GORM model for database
type EventModel struct {
	ID        uint      `gorm:"primaryKey"`
	Title     string    `gorm:"size:255;not null"`
	Content   string    `gorm:"type:text"`
	Cover     string    `gorm:"size:255"`
	Location  string    `gorm:"type:text;not null"`
	StartTime time.Time `gorm:"not null"`
	EndTime   time.Time `gorm:"not null"`
	UserID    uint      `gorm:"not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

func (EventModel) TableName() string {
	return "events"
}

type eventRepositoryImpl struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) repository.EventRepository {
	return &eventRepositoryImpl{db: db}
}

func (r *eventRepositoryImpl) Save(event *entity.Event) error {
	model := r.entityToModel(event)
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	event.ID = model.ID
	event.CreatedAt = model.CreatedAt
	event.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *eventRepositoryImpl) FindAll() ([]*entity.Event, error) {
	var models []EventModel
	if err := r.db.Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *eventRepositoryImpl) FindByID(id uint) (*entity.Event, error) {
	var model EventModel
	if err := r.db.First(&model, id).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *eventRepositoryImpl) FindByUserID(userID uint) ([]*entity.Event, error) {
	var models []EventModel
	if err := r.db.Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *eventRepositoryImpl) Update(event *entity.Event) error {
	model := r.entityToModel(event)
	return r.db.Save(&model).Error
}

func (r *eventRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&EventModel{}, id).Error
}

func (r *eventRepositoryImpl) entityToModel(event *entity.Event) *EventModel {
	return &EventModel{
		ID:        event.ID,
		Title:     event.Title,
		Content:   event.Content,
		Cover:     event.Cover,
		Location:  event.Location,
		StartTime: event.StartTime,
		EndTime:   event.EndTime,
		UserID:    event.UserID,
		CreatedAt: event.CreatedAt,
		UpdatedAt: event.UpdatedAt,
	}
}

func (r *eventRepositoryImpl) modelToEntity(model *EventModel) *entity.Event {
	return &entity.Event{
		ID:        model.ID,
		Title:     model.Title,
		Content:   model.Content,
		Cover:     model.Cover,
		Location:  model.Location,
		StartTime: model.StartTime,
		EndTime:   model.EndTime,
		UserID:    model.UserID,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
}

func (r *eventRepositoryImpl) modelsToEntities(models []EventModel) []*entity.Event {
	entities := make([]*entity.Event, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(&model)
	}
	return entities
}