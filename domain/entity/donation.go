package entity

import "time"

// DonationType represents the type of donation
type DonationType string

const (
	DonationTypeDonate  DonationType = "donate"
	DonationTypeSponsor DonationType = "sponsor"
)

// DonorType represents who is making the donation
type DonorType string

const (
	DonorTypeUser         DonorType = "user"
	DonorTypeCompany      DonorType = "company"
	DonorTypeOrganization DonorType = "organization"
	DonorTypeChurch       DonorType = "church"
)

// DonationStatus represents the status of a donation
type DonationStatus string

const (
	DonationStatusPending   DonationStatus = "pending"
	DonationStatusCompleted DonationStatus = "completed"
	DonationStatusFailed    DonationStatus = "failed"
	DonationStatusRefunded  DonationStatus = "refunded"
)

// Donation represents a donation or sponsorship
type Donation struct {
	ID             uint
	Type           DonationType   // donate or sponsor
	DonorType      DonorType      // user or company
	UserID         *uint          // For user donations
	SupporterID    *uint          // For company/organization donations (normalized)
	CompanyName    string         // For company/organization donations (legacy/inline)
	CompanyEmail   string         // For company/organization donations (legacy/inline)
	CompanyPhone   string         // For company/organization donations (legacy/inline)
	Amount         float64        // Donation amount
	Currency       string         // Currency code (e.g., USD, KHR)
	Message        string         // Optional message from donor
	Status         DonationStatus // Payment status
	TransactionID  string         // Payment transaction ID
	PaymentMethod  string         // Payment method used
	QRExpiresAt    *time.Time     // QR code expiration time (3 minutes from creation)
	EventID        *uint          // Optional: link to specific event
	User           *User          // Related user (if donor is user)
	Supporter      *Supporter     // Related supporter (if using normalized supporter)
	Event          *Event         // Related event (if donation is for specific event)
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// NewUserDonation creates a new donation from a user
func NewUserDonation(donationType DonationType, userID uint, amount float64, currency, message string) *Donation {
	return &Donation{
		Type:      donationType,
		DonorType: DonorTypeUser,
		UserID:    &userID,
		Amount:    amount,
		Currency:  currency,
		Message:   message,
		Status:    DonationStatusPending,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// NewCompanyDonation creates a new donation from a company/organization (inline details)
func NewCompanyDonation(donationType DonationType, companyName, companyEmail, companyPhone string, amount float64, currency, message string) *Donation {
	return &Donation{
		Type:         donationType,
		DonorType:    DonorTypeCompany,
		CompanyName:  companyName,
		CompanyEmail: companyEmail,
		CompanyPhone: companyPhone,
		Amount:       amount,
		Currency:     currency,
		Message:      message,
		Status:       DonationStatusPending,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// NewSupporterDonation creates a new donation from a supporter (normalized)
func NewSupporterDonation(donationType DonationType, supporterID uint, amount float64, currency, message string) *Donation {
	return &Donation{
		Type:        donationType,
		DonorType:   DonorTypeCompany,
		SupporterID: &supporterID,
		Amount:      amount,
		Currency:    currency,
		Message:     message,
		Status:      DonationStatusPending,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// IsValid validates donation entity
func (d *Donation) IsValid() bool {
	// Check donation type is valid
	if d.Type != DonationTypeDonate && d.Type != DonationTypeSponsor {
		return false
	}

	// Check donor type is valid
	if d.DonorType != DonorTypeUser && d.DonorType != DonorTypeCompany && d.DonorType != DonorTypeOrganization && d.DonorType != DonorTypeChurch {
		return false
	}

	// Check amount is positive
	if d.Amount <= 0 {
		return false
	}

	// Check currency is set
	if d.Currency == "" {
		return false
	}

	// For user donations, UserID must be set
	if d.DonorType == DonorTypeUser && (d.UserID == nil || *d.UserID == 0) {
		return false
	}

	// For company/organization/church donations, either supporter_id or (company name and email) must be set
	if d.DonorType == DonorTypeCompany || d.DonorType == DonorTypeOrganization || d.DonorType == DonorTypeChurch {
		hasSupporter := d.SupporterID != nil && *d.SupporterID > 0
		hasInlineInfo := d.CompanyName != "" && d.CompanyEmail != ""
		if !hasSupporter && !hasInlineInfo {
			return false
		}
	}

	return true
}

// Complete marks donation as completed
func (d *Donation) Complete(transactionID, paymentMethod string) {
	d.Status = DonationStatusCompleted
	d.TransactionID = transactionID
	d.PaymentMethod = paymentMethod
	d.UpdatedAt = time.Now()
}

// Fail marks donation as failed
func (d *Donation) Fail() {
	d.Status = DonationStatusFailed
	d.UpdatedAt = time.Now()
}

// Refund marks donation as refunded
func (d *Donation) Refund() {
	d.Status = DonationStatusRefunded
	d.UpdatedAt = time.Now()
}

// SetEvent links donation to a specific event
func (d *Donation) SetEvent(eventID uint) {
	d.EventID = &eventID
	d.UpdatedAt = time.Now()
}

// SetQRExpiration sets QR code expiration (3 minutes from now)
func (d *Donation) SetQRExpiration() {
	expiresAt := time.Now().Add(3 * time.Minute)
	d.QRExpiresAt = &expiresAt
	d.UpdatedAt = time.Now()
}

// IsQRExpired checks if QR code has expired
func (d *Donation) IsQRExpired() bool {
	if d.QRExpiresAt == nil {
		return false // No expiration set (card payment)
	}
	return time.Now().After(*d.QRExpiresAt)
}

// GetQRTimeRemaining returns remaining time before QR expires
func (d *Donation) GetQRTimeRemaining() time.Duration {
	if d.QRExpiresAt == nil {
		return 0
	}
	remaining := time.Until(*d.QRExpiresAt)
	if remaining < 0 {
		return 0
	}
	return remaining
}

