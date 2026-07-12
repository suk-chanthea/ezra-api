package entity

import "time"

type Church struct {
	ID              uint
	Name        string     // Church full name (unique)
	Address         string
	Phone           string
	Email           string
	Website         string
	PastorName      string     // Name of the pastor/priest
	Description     string     // About the church
	Logo            string     // Church logo/image URL
	EstablishedDate *time.Time // When the church was established
	Denomination    string     // Baptist, Catholic, Presbyterian, etc.
	OwnerID         *uint      // Church owner/admin who can approve members
	Owner           *User      // Church owner relationship
	Users           []User     // Members of this church (includes pending)
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func (c *Church) TableName() string {
	return "churches"
}

