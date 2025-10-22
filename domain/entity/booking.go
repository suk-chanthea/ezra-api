package entity

import "time"

// BookingStatus represents the status of a booking
type BookingStatus string

const (
	BookingStatusPending   BookingStatus = "pending"
	BookingStatusConfirmed BookingStatus = "confirmed"
	BookingStatusCancelled BookingStatus = "cancelled"
)

// Booking represents an event booking/registration
type Booking struct {
	ID        uint
	EventID   uint
	UserID    uint
	Status    BookingStatus
	Notes     string
	Event     *Event // Related event
	User      *User  // Related user
	CreatedAt time.Time
	UpdatedAt time.Time
}

// NewBooking creates a new booking entity
func NewBooking(eventID, userID uint, notes string) *Booking {
	return &Booking{
		EventID:   eventID,
		UserID:    userID,
		Status:    BookingStatusPending,
		Notes:     notes,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// IsValid validates booking entity
func (b *Booking) IsValid() bool {
	if b.EventID == 0 || b.UserID == 0 {
		return false
	}
	return true
}

// Confirm changes booking status to confirmed
func (b *Booking) Confirm() {
	b.Status = BookingStatusConfirmed
	b.UpdatedAt = time.Now()
}

// Cancel changes booking status to cancelled
func (b *Booking) Cancel() {
	b.Status = BookingStatusCancelled
	b.UpdatedAt = time.Now()
}

