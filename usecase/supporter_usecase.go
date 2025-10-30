package usecase

import (
	"errors"
	"fmt"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
)

type SupporterUseCase interface {
	CreateSupporter(req *dto.CreateSupporterRequest, userID *uint) (*dto.SupporterResponse, error)
	GetSupporterByID(id uint) (*dto.SupporterResponse, error)
	GetSupporterByEmail(email string) (*dto.SupporterResponse, error)
	GetAllSupporters(page, pageSize int) (*dto.PaginatedResponse, error)
	GetSupportersByType(supporterType string, page, pageSize int) (*dto.PaginatedResponse, error)
	GetSupportersByUser(userID uint, page, pageSize int) ([]*dto.SupporterResponse, error)
	UpdateSupporter(id uint, req *dto.UpdateSupporterRequest, userID *uint) (*dto.SupporterResponse, error)
	DeleteSupporter(id uint, userID *uint) error
	GetSupporterStats(supporterID uint) (map[string]interface{}, error)
}

type supporterUseCase struct {
	supporterRepo repository.SupporterRepository
	donationRepo  repository.DonationRepository
}

func NewSupporterUseCase(
	supporterRepo repository.SupporterRepository,
	donationRepo repository.DonationRepository,
) SupporterUseCase {
	return &supporterUseCase{
		supporterRepo: supporterRepo,
		donationRepo:  donationRepo,
	}
}

func (uc *supporterUseCase) CreateSupporter(req *dto.CreateSupporterRequest, userID *uint) (*dto.SupporterResponse, error) {
	// Check if supporter with this email already exists
	existing, _ := uc.supporterRepo.FindByEmail(req.Email)
	if existing != nil {
		return nil, errors.New("supporter with this email already exists")
	}

	supporter := &entity.Supporter{
		Name:        req.Name,
		Email:       req.Email,
		Phone:       req.Phone,
		Type:        entity.SupporterType(req.Type),
		Website:     req.Website,
		Address:     req.Address,
		Logo:        req.Logo,
		Description: req.Description,
		UserID:      userID,
	}

	if err := uc.supporterRepo.Create(supporter); err != nil {
		return nil, fmt.Errorf("failed to create supporter: %w", err)
	}

	return uc.entityToResponse(supporter), nil
}

func (uc *supporterUseCase) GetSupporterByID(id uint) (*dto.SupporterResponse, error) {
	supporter, err := uc.supporterRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("supporter not found")
	}

	response := uc.entityToResponse(supporter)

	// Get donation stats for this supporter
	stats, _ := uc.GetSupporterStats(id)
	if stats != nil {
		if totalDonations, ok := stats["total_donations"].(int64); ok {
			response.TotalDonations = int(totalDonations)
		}
		if totalAmount, ok := stats["total_amount"].(float64); ok {
			response.TotalAmount = totalAmount
		}
	}

	return response, nil
}

func (uc *supporterUseCase) GetSupporterByEmail(email string) (*dto.SupporterResponse, error) {
	supporter, err := uc.supporterRepo.FindByEmail(email)
	if err != nil {
		return nil, errors.New("supporter not found")
	}
	return uc.entityToResponse(supporter), nil
}

func (uc *supporterUseCase) GetAllSupporters(page, pageSize int) (*dto.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	supporters, err := uc.supporterRepo.FindAll(pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch supporters: %w", err)
	}

	responses := make([]interface{}, len(supporters))
	for i, supporter := range supporters {
		responses[i] = uc.entityToResponse(supporter)
	}

	totalCount, err := uc.supporterRepo.Count()
	if err != nil {
		return nil, err
	}

	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse{
		Data: responses,
		Pagination: &dto.PaginationMetadata{
			CurrentPage:  page,
			PageSize:     pageSize,
			TotalPages:   totalPages,
			TotalRecords: totalCount,
			HasNextPage:  page < totalPages,
			HasPrevPage:  page > 1,
		},
	}, nil
}

func (uc *supporterUseCase) GetSupportersByType(supporterType string, page, pageSize int) (*dto.PaginatedResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}

	offset := (page - 1) * pageSize
	supporters, err := uc.supporterRepo.FindByType(entity.SupporterType(supporterType), pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch supporters: %w", err)
	}

	responses := make([]interface{}, len(supporters))
	for i, supporter := range supporters {
		responses[i] = uc.entityToResponse(supporter)
	}

	totalCount, err := uc.supporterRepo.CountByType(entity.SupporterType(supporterType))
	if err != nil {
		return nil, err
	}

	totalPages := int(totalCount) / pageSize
	if int(totalCount)%pageSize > 0 {
		totalPages++
	}

	return &dto.PaginatedResponse{
		Data: responses,
		Pagination: &dto.PaginationMetadata{
			CurrentPage:  page,
			PageSize:     pageSize,
			TotalPages:   totalPages,
			TotalRecords: totalCount,
			HasNextPage:  page < totalPages,
			HasPrevPage:  page > 1,
		},
	}, nil
}

func (uc *supporterUseCase) GetSupportersByUser(userID uint, page, pageSize int) ([]*dto.SupporterResponse, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 20
	}

	offset := (page - 1) * pageSize
	supporters, err := uc.supporterRepo.FindByUser(userID, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch supporters: %w", err)
	}

	responses := make([]*dto.SupporterResponse, len(supporters))
	for i, supporter := range supporters {
		responses[i] = uc.entityToResponse(supporter)
	}

	return responses, nil
}

func (uc *supporterUseCase) UpdateSupporter(id uint, req *dto.UpdateSupporterRequest, userID *uint) (*dto.SupporterResponse, error) {
	supporter, err := uc.supporterRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("supporter not found")
	}

	// Check if user has permission to update (must be the owner)
	if userID != nil && supporter.UserID != nil && *supporter.UserID != *userID {
		return nil, errors.New("you don't have permission to update this supporter")
	}

	// Check if email is being changed and if it already exists
	if req.Email != supporter.Email {
		existing, _ := uc.supporterRepo.FindByEmail(req.Email)
		if existing != nil && existing.ID != id {
			return nil, errors.New("supporter with this email already exists")
		}
	}

	supporter.Name = req.Name
	supporter.Email = req.Email
	supporter.Phone = req.Phone
	supporter.Type = entity.SupporterType(req.Type)
	supporter.Website = req.Website
	supporter.Address = req.Address
	supporter.Logo = req.Logo
	supporter.Description = req.Description

	if err := uc.supporterRepo.Update(supporter); err != nil {
		return nil, fmt.Errorf("failed to update supporter: %w", err)
	}

	return uc.entityToResponse(supporter), nil
}

func (uc *supporterUseCase) DeleteSupporter(id uint, userID *uint) error {
	supporter, err := uc.supporterRepo.FindByID(id)
	if err != nil {
		return errors.New("supporter not found")
	}

	// Check if user has permission to delete (must be the owner)
	if userID != nil && supporter.UserID != nil && *supporter.UserID != *userID {
		return errors.New("you don't have permission to delete this supporter")
	}

	return uc.supporterRepo.Delete(id)
}

func (uc *supporterUseCase) GetSupporterStats(supporterID uint) (map[string]interface{}, error) {
	// Count total donations from this supporter
	donations, err := uc.donationRepo.FindAll(0, 0)
	if err != nil {
		return nil, err
	}

	var totalDonations int64
	var totalAmount float64

	for _, donation := range donations {
		if donation.SupporterID != nil && *donation.SupporterID == supporterID && donation.Status == entity.DonationStatusCompleted {
			totalDonations++
			totalAmount += donation.Amount
		}
	}

	return map[string]interface{}{
		"total_donations": totalDonations,
		"total_amount":    totalAmount,
	}, nil
}

func (uc *supporterUseCase) entityToResponse(supporter *entity.Supporter) *dto.SupporterResponse {
	response := &dto.SupporterResponse{
		ID:          supporter.ID,
		Name:        supporter.Name,
		Email:       supporter.Email,
		Phone:       supporter.Phone,
		Type:        string(supporter.Type),
		Website:     supporter.Website,
		Address:     supporter.Address,
		Logo:        supporter.Logo,
		Description: supporter.Description,
		UserID:      supporter.UserID,
		CreatedAt:   supporter.CreatedAt,
		UpdatedAt:   supporter.UpdatedAt,
	}

	// Convert user if present
	if supporter.User != nil {
		response.User = &dto.UserResponse{
			ID:       supporter.User.ID,
			Username: supporter.User.Username,
			Fullname: supporter.User.Fullname,
			Email:    supporter.User.Email,
			Profile:  supporter.User.Profile,
			Role:     supporter.User.Role,
		}
	}

	return response
}

