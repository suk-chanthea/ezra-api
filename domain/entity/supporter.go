package entity

import "time"

type SupporterType string

const (
	SupporterTypeCompany      SupporterType = "company"
	SupporterTypeOrganization SupporterType = "organization"
	SupporterTypeChurch       SupporterType = "church"
)

type Supporter struct {
	ID          uint
	Name        string // Company/Organization name
	Email       string
	Phone       string
	Type        SupporterType
	Website     string
	Address     string
	Logo        string
	Description string
	UserID      *uint // Optional: user who manages this supporter
	User        *User // Related user
	Donations   []Donation // All donations from this supporter
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (s *Supporter) TableName() string {
	return "supporters"
}

