package usecase

import (
	"errors"
	"fmt"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
	"github.com/suk-chanthea/ezra/infrastructure/payment"
)

type DonationUseCase interface {
	CreateDonation(req *dto.CreateDonationRequest, userID *uint) (*dto.DonationResponse, error)
	InitiatePayment(donationID uint) (*dto.InitiatePaymentResponse, error)
	HandlePaymentCallback(transactionID, status, approvalCode, paymentMethod string) error
	GetAllDonations(filter *dto.DonationFilterRequest) ([]*dto.DonationResponse, *dto.PaginationMetadata, error)
	GetDonationByID(id uint) (*dto.DonationResponse, error)
	GetDonationsByUserID(userID uint, limit, offset int) ([]*dto.DonationResponse, error)
	GetDonationsByType(donationType string, limit, offset int) ([]*dto.DonationResponse, error)
	GetDonationsByEventID(eventID uint, limit, offset int) ([]*dto.DonationResponse, error)
	UpdateDonationStatus(id uint, req *dto.UpdateDonationStatusRequest) error
	DeleteDonation(id uint, userID uint) error
	GetDonationStats() (*dto.DonationStatsResponse, error)
	GetDonationStatsByEventID(eventID uint) (*dto.DonationStatsResponse, error)
}

type donationUseCase struct {
	donationRepo   repository.DonationRepository
	userRepo       repository.UserRepository
	eventRepo      repository.EventRepository
	paywayService  payment.PaywayService
}

func NewDonationUseCase(
	donationRepo repository.DonationRepository,
	userRepo repository.UserRepository,
	eventRepo repository.EventRepository,
	paywayService payment.PaywayService,
) DonationUseCase {
	return &donationUseCase{
		donationRepo:  donationRepo,
		userRepo:      userRepo,
		eventRepo:     eventRepo,
		paywayService: paywayService,
	}
}

func (uc *donationUseCase) CreateDonation(req *dto.CreateDonationRequest, userID *uint) (*dto.DonationResponse, error) {
	var donation *entity.Donation

	donationType := entity.DonationType(req.Type)
	
	// Create donation based on donor type
	if req.DonorType == "user" {
		// For user donations, userID must be provided
		if userID == nil || *userID == 0 {
			return nil, errors.New("user must be authenticated to make a donation")
		}

		// Verify user exists
		_, err := uc.userRepo.FindByID(*userID)
		if err != nil {
			return nil, errors.New("user not found")
		}

		donation = entity.NewUserDonation(donationType, *userID, req.Amount, req.Currency, req.Message)
	} else if req.DonorType == "company" || req.DonorType == "organization" || req.DonorType == "church" {
		// For company/organization/church donations, company info must be provided
		if req.CompanyName == "" || req.CompanyEmail == "" {
			return nil, errors.New("company/organization/church name and email are required")
		}

		donation = entity.NewCompanyDonation(donationType, req.CompanyName, req.CompanyEmail, req.CompanyPhone, req.Amount, req.Currency, req.Message)
		donation.DonorType = entity.DonorType(req.DonorType) // Set the correct donor type
	} else {
		return nil, errors.New("invalid donor type")
	}

	// Link to event if provided
	if req.EventID != nil && *req.EventID > 0 {
		// Verify event exists
		_, err := uc.eventRepo.FindByID(*req.EventID)
		if err != nil {
			return nil, errors.New("event not found")
		}
		donation.SetEvent(*req.EventID)
	}

	// Validate
	if !donation.IsValid() {
		return nil, errors.New("invalid donation data")
	}

	// Save
	if err := uc.donationRepo.Save(donation); err != nil {
		return nil, err
	}

	// Fetch the created donation with relations
	createdDonation, err := uc.donationRepo.FindByID(donation.ID)
	if err != nil {
		return nil, err
	}

	response := uc.entityToResponse(createdDonation)

	// If initiate_payment is true, automatically initiate payment
	if req.InitiatePayment {
		paymentInfo, err := uc.InitiatePayment(createdDonation.ID)
		if err != nil {
			// Don't fail the whole request if payment initiation fails
			// Just log the error and return without payment info
			fmt.Printf("Failed to initiate payment: %v\n", err)
		} else {
			response.PaymentInfo = paymentInfo
		}
	}

	return response, nil
}

func (uc *donationUseCase) InitiatePayment(donationID uint) (*dto.InitiatePaymentResponse, error) {
	// Get donation
	donation, err := uc.donationRepo.FindByID(donationID)
	if err != nil {
		return nil, errors.New("donation not found")
	}

	// Check if payment is already completed
	if donation.Status == entity.DonationStatusCompleted {
		return nil, errors.New("donation already paid")
	}

	// Generate transaction ID
	transactionID := fmt.Sprintf("DON-%d-%d", donation.ID, donation.CreatedAt.Unix())

	// Get customer info
	var customerName, customerEmail, customerPhone string
	if donation.DonorType == entity.DonorTypeUser && donation.User != nil {
		customerName = donation.User.Fullname
		customerEmail = donation.User.Email
		customerPhone = "" // User might not have phone
	} else {
		customerName = donation.CompanyName
		customerEmail = donation.CompanyEmail
		customerPhone = donation.CompanyPhone
	}

	// Prepare items description
	items := fmt.Sprintf("%s - Donation #%d", donation.Type, donation.ID)
	if donation.Event != nil {
		items = fmt.Sprintf("%s for %s", donation.Type, donation.Event.Title)
	}

	// Format amount
	amountStr := payment.FormatAmount(donation.Amount, donation.Currency)

	var paymentResp *payment.PaymentResponse

	// Initiate payment based on donation type
	if donation.Type == entity.DonationTypeDonate {
		// Donate uses QR code with 3-minute expiration
		paymentResp, err = uc.paywayService.InitiateQRPayment(
			transactionID,
			amountStr,
			donation.Currency,
			customerName,
			customerEmail,
			customerPhone,
			items,
		)
		
		// Set QR expiration (3 minutes)
		if err == nil {
			donation.SetQRExpiration()
		}
	} else if donation.Type == entity.DonationTypeSponsor {
		// Sponsor uses card payment (no expiration needed)
		paymentResp, err = uc.paywayService.InitiateCardPayment(
			transactionID,
			amountStr,
			donation.Currency,
			customerName,
			customerEmail,
			customerPhone,
			items,
		)
	} else {
		return nil, errors.New("invalid donation type")
	}

	if err != nil {
		return nil, fmt.Errorf("failed to initiate payment: %w", err)
	}

	// Update donation with transaction ID and expiration
	donation.TransactionID = transactionID
	if err := uc.donationRepo.Update(donation); err != nil {
		return nil, fmt.Errorf("failed to update donation: %w", err)
	}

	// Build response
	response := &dto.InitiatePaymentResponse{
		DonationID:    donation.ID,
		TransactionID: transactionID,
		Amount:        amountStr,
		Currency:      donation.Currency,
		Message:       "Payment initiated successfully",
	}

	if donation.Type == entity.DonationTypeDonate {
		response.PaymentMethod = "qr"
		response.QRCode = paymentResp.QRCode
		response.ExpiresAt = donation.QRExpiresAt
		if donation.QRExpiresAt != nil {
			response.ExpiresInSeconds = int(donation.GetQRTimeRemaining().Seconds())
		}
	} else {
		response.PaymentMethod = "card"
		response.PaymentURL = paymentResp.PaymentURL
	}

	return response, nil
}

func (uc *donationUseCase) HandlePaymentCallback(transactionID, status, approvalCode, paymentMethod string) error {
	// Find donation by transaction ID
	donations, err := uc.donationRepo.FindAll(-1, 0)
	if err != nil {
		return err
	}

	var donation *entity.Donation
	for _, d := range donations {
		if d.TransactionID == transactionID {
			donation = d
			break
		}
	}

	if donation == nil {
		return errors.New("donation not found for transaction")
	}

	// Update status based on payment result
	if status == "success" {
		donation.Complete(approvalCode, paymentMethod)
	} else {
		donation.Fail()
	}

	return uc.donationRepo.Update(donation)
}

func (uc *donationUseCase) GetAllDonations(filter *dto.DonationFilterRequest) ([]*dto.DonationResponse, *dto.PaginationMetadata, error) {
	var donations []*entity.Donation
	var err error

	limit := filter.GetPageSize()
	offset := filter.GetOffset()

	// Apply filters
	if filter.Type != "" {
		donations, err = uc.donationRepo.FindByType(entity.DonationType(filter.Type), limit, offset)
	} else if filter.DonorType != "" {
		donations, err = uc.donationRepo.FindByDonorType(entity.DonorType(filter.DonorType), limit, offset)
	} else if filter.Status != "" {
		donations, err = uc.donationRepo.FindByStatus(entity.DonationStatus(filter.Status), limit, offset)
	} else if filter.EventID != nil && *filter.EventID > 0 {
		donations, err = uc.donationRepo.FindByEventID(*filter.EventID, limit, offset)
	} else {
		donations, err = uc.donationRepo.FindAll(limit, offset)
	}

	if err != nil {
		return nil, nil, err
	}

	// Get total count for pagination
	total, err := uc.donationRepo.Count()
	if err != nil {
		return nil, nil, err
	}

	pagination := dto.NewPaginationMetadata(filter.GetPage(), filter.GetPageSize(), total)
	return uc.entitiesToResponses(donations), pagination, nil
}

func (uc *donationUseCase) GetDonationByID(id uint) (*dto.DonationResponse, error) {
	donation, err := uc.donationRepo.FindByID(id)
	if err != nil {
		return nil, err
	}

	return uc.entityToResponse(donation), nil
}

func (uc *donationUseCase) GetDonationsByUserID(userID uint, limit, offset int) ([]*dto.DonationResponse, error) {
	donations, err := uc.donationRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	return uc.entitiesToResponses(donations), nil
}

func (uc *donationUseCase) GetDonationsByType(donationType string, limit, offset int) ([]*dto.DonationResponse, error) {
	donations, err := uc.donationRepo.FindByType(entity.DonationType(donationType), limit, offset)
	if err != nil {
		return nil, err
	}

	return uc.entitiesToResponses(donations), nil
}

func (uc *donationUseCase) GetDonationsByEventID(eventID uint, limit, offset int) ([]*dto.DonationResponse, error) {
	donations, err := uc.donationRepo.FindByEventID(eventID, limit, offset)
	if err != nil {
		return nil, err
	}

	return uc.entitiesToResponses(donations), nil
}

func (uc *donationUseCase) UpdateDonationStatus(id uint, req *dto.UpdateDonationStatusRequest) error {
	// Check if donation exists
	donation, err := uc.donationRepo.FindByID(id)
	if err != nil {
		return errors.New("donation not found")
	}

	// Update status
	status := entity.DonationStatus(req.Status)
	
	// Use entity methods for specific status changes
	switch status {
	case entity.DonationStatusCompleted:
		donation.Complete(req.TransactionID, req.PaymentMethod)
	case entity.DonationStatusFailed:
		donation.Fail()
	case entity.DonationStatusRefunded:
		donation.Refund()
	default:
		donation.Status = status
	}

	// Save
	return uc.donationRepo.Update(donation)
}

func (uc *donationUseCase) DeleteDonation(id uint, userID uint) error {
	// Check if donation exists
	donation, err := uc.donationRepo.FindByID(id)
	if err != nil {
		return errors.New("donation not found")
	}

	// Only allow deletion of user's own donations or by admins
	// For company donations, only admins can delete
	if donation.DonorType == entity.DonorTypeUser {
		if donation.UserID == nil || *donation.UserID != userID {
			return errors.New("unauthorized to delete this donation")
		}
	}

	return uc.donationRepo.Delete(id)
}

func (uc *donationUseCase) GetDonationStats() (*dto.DonationStatsResponse, error) {
	totalAmount, err := uc.donationRepo.GetTotalAmount()
	if err != nil {
		return nil, err
	}

	donateAmount, err := uc.donationRepo.GetTotalAmountByType(entity.DonationTypeDonate)
	if err != nil {
		return nil, err
	}

	sponsorAmount, err := uc.donationRepo.GetTotalAmountByType(entity.DonationTypeSponsor)
	if err != nil {
		return nil, err
	}

	donateCount, err := uc.donationRepo.CountByType(entity.DonationTypeDonate)
	if err != nil {
		return nil, err
	}

	sponsorCount, err := uc.donationRepo.CountByType(entity.DonationTypeSponsor)
	if err != nil {
		return nil, err
	}

	// Count by donor type (we'll need to fetch and count manually)
	userDonations, err := uc.donationRepo.FindByDonorType(entity.DonorTypeUser, -1, 0)
	if err != nil {
		return nil, err
	}

	companyDonations, err := uc.donationRepo.FindByDonorType(entity.DonorTypeCompany, -1, 0)
	if err != nil {
		return nil, err
	}

	return &dto.DonationStatsResponse{
		TotalAmount:      totalAmount,
		TotalDonations:   donateCount,
		TotalSponsors:    sponsorCount,
		DonateAmount:     donateAmount,
		SponsorAmount:    sponsorAmount,
		UserDonations:    int64(len(userDonations)),
		CompanyDonations: int64(len(companyDonations)),
	}, nil
}

func (uc *donationUseCase) GetDonationStatsByEventID(eventID uint) (*dto.DonationStatsResponse, error) {
	totalAmount, err := uc.donationRepo.GetTotalAmountByEventID(eventID)
	if err != nil {
		return nil, err
	}

	// Get all donations for the event
	donations, err := uc.donationRepo.FindByEventID(eventID, -1, 0)
	if err != nil {
		return nil, err
	}

	// Calculate stats
	var donateAmount, sponsorAmount float64
	var donateCount, sponsorCount, userCount, companyCount int64

	for _, donation := range donations {
		if donation.Type == entity.DonationTypeDonate {
			donateAmount += donation.Amount
			donateCount++
		} else if donation.Type == entity.DonationTypeSponsor {
			sponsorAmount += donation.Amount
			sponsorCount++
		}

		if donation.DonorType == entity.DonorTypeUser {
			userCount++
		} else if donation.DonorType == entity.DonorTypeCompany {
			companyCount++
		}
	}

	return &dto.DonationStatsResponse{
		TotalAmount:      totalAmount,
		TotalDonations:   donateCount,
		TotalSponsors:    sponsorCount,
		DonateAmount:     donateAmount,
		SponsorAmount:    sponsorAmount,
		UserDonations:    userCount,
		CompanyDonations: companyCount,
	}, nil
}

func (uc *donationUseCase) entityToResponse(donation *entity.Donation) *dto.DonationResponse {
	response := &dto.DonationResponse{
		ID:            donation.ID,
		Type:          string(donation.Type),
		DonorType:     string(donation.DonorType),
		UserID:        donation.UserID,
		CompanyName:   donation.CompanyName,
		CompanyEmail:  donation.CompanyEmail,
		CompanyPhone:  donation.CompanyPhone,
		Amount:        donation.Amount,
		Currency:      donation.Currency,
		Message:       donation.Message,
		Status:        string(donation.Status),
		TransactionID: donation.TransactionID,
		PaymentMethod: donation.PaymentMethod,
		EventID:       donation.EventID,
		CreatedAt:     donation.CreatedAt,
		UpdatedAt:     donation.UpdatedAt,
	}

	// Convert user if available
	if donation.User != nil {
		response.User = &dto.UserResponse{
			ID:        donation.User.ID,
			Username:  donation.User.Username,
			Fullname:  donation.User.Fullname,
			Profile:   donation.User.Profile,
			Email:     donation.User.Email,
			Role:      donation.User.Role,
			CreatedAt: donation.User.CreatedAt,
			UpdatedAt: donation.User.UpdatedAt,
		}
	}

	// Convert event if available
	if donation.Event != nil {
		response.Event = &dto.EventResponse{
			ID:        donation.Event.ID,
			Title:     donation.Event.Title,
			Content:   donation.Event.Content,
			Cover:     donation.Event.Cover,
			Location:  donation.Event.Location,
			StartTime: donation.Event.StartTime,
			EndTime:   donation.Event.EndTime,
			UserID:    donation.Event.UserID,
			CreatedAt: donation.Event.CreatedAt,
			UpdatedAt: donation.Event.UpdatedAt,
		}
	}

	return response
}

func (uc *donationUseCase) entitiesToResponses(donations []*entity.Donation) []*dto.DonationResponse {
	responses := make([]*dto.DonationResponse, len(donations))
	for i, donation := range donations {
		responses[i] = uc.entityToResponse(donation)
	}
	return responses
}

