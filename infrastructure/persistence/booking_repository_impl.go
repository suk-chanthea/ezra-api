package persistence

import (
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

// BookingModel is the GORM model for database
type BookingModel struct {
	ID        uint      `gorm:"primaryKey"`
	EventID   uint      `gorm:"not null;index"`
	UserID    uint      `gorm:"not null;index"`
	Status    string    `gorm:"size:50;not null;default:'pending'"`
	Notes     string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
	
	// Relations
	Event EventModel `gorm:"foreignKey:EventID"`
	User  UserModel  `gorm:"foreignKey:UserID"`
}

func (BookingModel) TableName() string {
	return "bookings"
}

type bookingRepositoryImpl struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) repository.BookingRepository {
	return &bookingRepositoryImpl{db: db}
}

func (r *bookingRepositoryImpl) Save(booking *entity.Booking) error {
	model := r.entityToModel(booking)
	if err := r.db.Create(&model).Error; err != nil {
		return err
	}
	
	booking.ID = model.ID
	booking.CreatedAt = model.CreatedAt
	booking.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *bookingRepositoryImpl) FindAll() ([]*entity.Booking, error) {
	var models []BookingModel
	if err := r.db.Preload("Event").Preload("Event.Musics").Preload("User").Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *bookingRepositoryImpl) FindAllPaginated(offset, limit int) ([]*entity.Booking, int64, error) {
	var models []BookingModel
	var total int64
	
	// Get total count
	if err := r.db.Model(&BookingModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results with preloaded relations
	if err := r.db.Preload("Event").Preload("Event.Musics").Preload("User").Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	
	return r.modelsToEntities(models), total, nil
}

func (r *bookingRepositoryImpl) FindByID(id uint) (*entity.Booking, error) {
	var model BookingModel
	if err := r.db.Preload("Event").Preload("Event.Musics").Preload("User").First(&model, id).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *bookingRepositoryImpl) FindByEventID(eventID uint) ([]*entity.Booking, error) {
	var models []BookingModel
	if err := r.db.Preload("Event").Preload("Event.Musics").Preload("User").
		Where("event_id = ?", eventID).Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *bookingRepositoryImpl) FindByUserID(userID uint) ([]*entity.Booking, error) {
	var models []BookingModel
	if err := r.db.Preload("Event").Preload("Event.Musics").Preload("User").
		Where("user_id = ?", userID).Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *bookingRepositoryImpl) FindByEventAndUser(eventID, userID uint) (*entity.Booking, error) {
	var model BookingModel
	if err := r.db.Preload("Event").Preload("Event.Musics").Preload("User").
		Where("event_id = ? AND user_id = ?", eventID, userID).First(&model).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *bookingRepositoryImpl) Update(booking *entity.Booking) error {
	model := r.entityToModel(booking)
	return r.db.Save(&model).Error
}

func (r *bookingRepositoryImpl) Delete(id uint) error {
	return r.db.Delete(&BookingModel{}, id).Error
}

func (r *bookingRepositoryImpl) entityToModel(booking *entity.Booking) *BookingModel {
	return &BookingModel{
		ID:        booking.ID,
		EventID:   booking.EventID,
		UserID:    booking.UserID,
		Status:    string(booking.Status),
		Notes:     booking.Notes,
		CreatedAt: booking.CreatedAt,
		UpdatedAt: booking.UpdatedAt,
	}
}

func (r *bookingRepositoryImpl) modelToEntity(model *BookingModel) *entity.Booking {
	booking := &entity.Booking{
		ID:        model.ID,
		EventID:   model.EventID,
		UserID:    model.UserID,
		Status:    entity.BookingStatus(model.Status),
		Notes:     model.Notes,
		CreatedAt: model.CreatedAt,
		UpdatedAt: model.UpdatedAt,
	}
	
	// Convert event if loaded
	if model.Event.ID != 0 {
		booking.Event = &entity.Event{
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
		
		// Convert musics if loaded
		if len(model.Event.Musics) > 0 {
			booking.Event.Musics = make([]*entity.Music, len(model.Event.Musics))
			for i, musicModel := range model.Event.Musics {
				booking.Event.Musics[i] = &entity.Music{
					ID:          musicModel.ID,
					Title:       musicModel.Title,
					Artist:      musicModel.Artist,
					Album:       musicModel.Album,
					Genre:       musicModel.Genre,
					Duration:    musicModel.Duration,
					BPM:         musicModel.BPM,
					Key:         musicModel.Key,
					Cover:       musicModel.Cover,
					Lyrics:      musicModel.Lyrics,
					Description: musicModel.Description,
					UserID:      musicModel.UserID,
					CreatedAt:   musicModel.CreatedAt,
					UpdatedAt:   musicModel.UpdatedAt,
				}
			}
		}
	}
	
	// Convert user if loaded
	if model.User.ID != 0 {
		booking.User = &entity.User{
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
	
	return booking
}

func (r *bookingRepositoryImpl) modelsToEntities(models []BookingModel) []*entity.Booking {
	entities := make([]*entity.Booking, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(&model)
	}
	return entities
}

