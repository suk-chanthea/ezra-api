package persistence

import (
	"time"

	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"gorm.io/gorm"
)

// EventModel is the GORM model for database
type EventModel struct {
	ID        uint         `gorm:"primaryKey"`
	Title     string       `gorm:"size:255;not null"`
	Content   string       `gorm:"type:text"`
	Cover     string       `gorm:"size:255"`
	Location  string       `gorm:"type:text;not null"`
	StartTime time.Time    `gorm:"not null"`
	EndTime   time.Time    `gorm:"not null"`
	UserID    uint         `gorm:"not null"`
	Musics    []MusicModel `gorm:"many2many:event_musics;joinForeignKey:EventID;joinReferences:MusicID"`
	CreatedAt time.Time    `gorm:"autoCreateTime"`
	UpdatedAt time.Time    `gorm:"autoUpdateTime"`
}

func (EventModel) TableName() string {
	return "events"
}

// EventMusicModel represents the junction table
type EventMusicModel struct {
	ID           uint      `gorm:"primaryKey"`
	EventID      uint      `gorm:"not null"`
	MusicID      uint      `gorm:"not null"`
	DisplayOrder int       `gorm:"default:0"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

func (EventMusicModel) TableName() string {
	return "event_musics"
}

type eventRepositoryImpl struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) repository.EventRepository {
	return &eventRepositoryImpl{db: db}
}

func (r *eventRepositoryImpl) Save(event *entity.Event) error {
	model := r.entityToModel(event)
	
	// Start transaction
	tx := r.db.Begin()
	
	// Create event
	if err := tx.Create(&model).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// Add musics if provided
	if len(event.MusicIDs) > 0 {
		if err := r.addMusicsToEventTx(tx, model.ID, event.MusicIDs); err != nil {
			tx.Rollback()
			return err
		}
	}
	
	tx.Commit()
	
	event.ID = model.ID
	event.CreatedAt = model.CreatedAt
	event.UpdatedAt = model.UpdatedAt
	return nil
}

func (r *eventRepositoryImpl) FindAll() ([]*entity.Event, error) {
	var models []EventModel
	if err := r.db.Preload("Musics").Find(&models).Error; err != nil {
		return nil, err
	}
	return r.modelsToEntities(models), nil
}

func (r *eventRepositoryImpl) FindAllPaginated(offset, limit int) ([]*entity.Event, int64, error) {
	var models []EventModel
	var total int64
	
	// Get total count
	if err := r.db.Model(&EventModel{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}
	
	// Get paginated results with preloaded musics
	if err := r.db.Preload("Musics").Offset(offset).Limit(limit).Find(&models).Error; err != nil {
		return nil, 0, err
	}
	
	return r.modelsToEntities(models), total, nil
}

func (r *eventRepositoryImpl) FindByID(id uint) (*entity.Event, error) {
	var model EventModel
	if err := r.db.Preload("Musics").First(&model, id).Error; err != nil {
		return nil, err
	}
	return r.modelToEntity(&model), nil
}

func (r *eventRepositoryImpl) FindByUserID(userID uint) ([]*entity.Event, error) {
	var models []EventModel
	if err := r.db.Preload("Musics").Where("user_id = ?", userID).Find(&models).Error; err != nil {
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

func (r *eventRepositoryImpl) AddMusicsToEvent(eventID uint, musicIDs []uint) error {
	return r.addMusicsToEventTx(r.db, eventID, musicIDs)
}

func (r *eventRepositoryImpl) addMusicsToEventTx(tx *gorm.DB, eventID uint, musicIDs []uint) error {
	for i, musicID := range musicIDs {
		eventMusic := EventMusicModel{
			EventID:      eventID,
			MusicID:      musicID,
			DisplayOrder: i,
		}
		if err := tx.Create(&eventMusic).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *eventRepositoryImpl) RemoveMusicsFromEvent(eventID uint, musicIDs []uint) error {
	return r.db.Where("event_id = ? AND music_id IN ?", eventID, musicIDs).
		Delete(&EventMusicModel{}).Error
}

func (r *eventRepositoryImpl) GetEventMusics(eventID uint) ([]*entity.Music, error) {
	var musicModels []MusicModel
	err := r.db.Table("musics").
		Joins("JOIN event_musics ON musics.id = event_musics.music_id").
		Where("event_musics.event_id = ?", eventID).
		Order("event_musics.display_order").
		Find(&musicModels).Error
	
	if err != nil {
		return nil, err
	}
	
	musics := make([]*entity.Music, len(musicModels))
	for i, model := range musicModels {
		musics[i] = &entity.Music{
			ID:          model.ID,
			Title:       model.Title,
			Artist:      model.Artist,
			Album:       model.Album,
			Genre:       model.Genre,
			Duration:    model.Duration,
			BPM:         model.BPM,
			Key:         model.Key,
			Cover:       model.Cover,
			Lyrics:      model.Lyrics,
			Description: model.Description,
			UserID:      model.UserID,
			CreatedAt:   model.CreatedAt,
			UpdatedAt:   model.UpdatedAt,
		}
	}
	return musics, nil
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
	event := &entity.Event{
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
	
	// Convert musics
	if len(model.Musics) > 0 {
		event.Musics = make([]*entity.Music, len(model.Musics))
		for i, musicModel := range model.Musics {
			event.Musics[i] = &entity.Music{
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
	
	return event
}

func (r *eventRepositoryImpl) modelsToEntities(models []EventModel) []*entity.Event {
	entities := make([]*entity.Event, len(models))
	for i, model := range models {
		entities[i] = r.modelToEntity(&model)
	}
	return entities
}