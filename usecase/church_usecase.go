package usecase

import (
	"errors"
	"fmt"
	"time"

	"github.com/suk-chanthea/ezra/domain/dto"
	"github.com/suk-chanthea/ezra/domain/entity"
	"github.com/suk-chanthea/ezra/domain/repository"
)

type ChurchUseCase interface {
	CreateChurch(req *dto.CreateChurchRequest, ownerID uint) (*dto.ChurchResponse, error)
	GetChurchByID(id uint) (*dto.ChurchResponse, error)
	GetAllChurches(page, pageSize int) (*dto.PaginatedResponse, error)
	GetChurchesByDenomination(denomination string, page, pageSize int) (*dto.PaginatedResponse, error)
	UpdateChurch(id uint, req *dto.UpdateChurchRequest, userID uint) (*dto.ChurchResponse, error)
	DeleteChurch(id uint, userID uint) error
	JoinChurch(userID uint, churchID uint) error
	LeaveChurch(userID uint) error
	ApproveMember(churchID uint, userID uint, targetUserID uint, status string) error
	GetPendingMembers(churchID uint, ownerID uint, page, pageSize int) (*dto.PaginatedResponse, error)
	GetApprovedMembers(churchID uint, page, pageSize int) (*dto.PaginatedResponse, error)
}

type churchUseCase struct {
	churchRepo repository.ChurchRepository
	userRepo   repository.UserRepository
}

func NewChurchUseCase(churchRepo repository.ChurchRepository, userRepo repository.UserRepository) ChurchUseCase {
	return &churchUseCase{
		churchRepo: churchRepo,
		userRepo:   userRepo,
	}
}

func (uc *churchUseCase) CreateChurch(req *dto.CreateChurchRequest, ownerID uint) (*dto.ChurchResponse, error) {
	// Check if church with this name already exists
	existing, _ := uc.churchRepo.FindByName(req.Name)
	if existing != nil {
		return nil, errors.New("church with this name already exists")
	}

	church := &entity.Church{
		Name:     req.Name,
		Address:      req.Address,
		Phone:        req.Phone,
		Email:        req.Email,
		Website:      req.Website,
		PastorName:   req.PastorName,
		Description:  req.Description,
		Logo:         req.Logo,
		Denomination: req.Denomination,
		OwnerID:      &ownerID,
	}

	// Parse established date if provided
	if req.EstablishedDate != "" {
		establishedDate, err := time.Parse("2006-01-02", req.EstablishedDate)
		if err == nil {
			church.EstablishedDate = &establishedDate
		}
	}

	if err := uc.churchRepo.Create(church); err != nil {
		return nil, fmt.Errorf("failed to create church: %w", err)
	}

	return uc.entityToResponse(church), nil
}

func (uc *churchUseCase) GetChurchByID(id uint) (*dto.ChurchResponse, error) {
	church, err := uc.churchRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("church not found")
	}

	return uc.entityToResponse(church), nil
}

func (uc *churchUseCase) GetAllChurches(page, pageSize int) (*dto.PaginatedResponse, error) {
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
	churches, err := uc.churchRepo.FindAll(pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch churches: %w", err)
	}

	responses := make([]interface{}, len(churches))
	for i, church := range churches {
		responses[i] = uc.entityToResponse(church)
	}

	totalCount, err := uc.churchRepo.Count()
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

func (uc *churchUseCase) GetChurchesByDenomination(denomination string, page, pageSize int) (*dto.PaginatedResponse, error) {
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
	churches, err := uc.churchRepo.FindByDenomination(denomination, pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch churches: %w", err)
	}

	responses := make([]interface{}, len(churches))
	for i, church := range churches {
		responses[i] = uc.entityToResponse(church)
	}

	// For simplicity, we don't have count by denomination, so we'll use total count
	totalCount := int64(len(churches))

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

func (uc *churchUseCase) UpdateChurch(id uint, req *dto.UpdateChurchRequest, userID uint) (*dto.ChurchResponse, error) {
	church, err := uc.churchRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("church not found")
	}

	// Check if user is the owner
	if church.OwnerID == nil || *church.OwnerID != userID {
		return nil, errors.New("only church owner can update church details")
	}

	// Check if name is being changed and if it already exists
	if req.Name != church.Name {
		existing, _ := uc.churchRepo.FindByName(req.Name)
		if existing != nil && existing.ID != id {
			return nil, errors.New("church with this name already exists")
		}
	}

	church.Name = req.Name
	church.Address = req.Address
	church.Phone = req.Phone
	church.Email = req.Email
	church.Website = req.Website
	church.PastorName = req.PastorName
	church.Description = req.Description
	church.Logo = req.Logo
	church.Denomination = req.Denomination

	// Parse established date if provided
	if req.EstablishedDate != "" {
		establishedDate, err := time.Parse("2006-01-02", req.EstablishedDate)
		if err == nil {
			church.EstablishedDate = &establishedDate
		}
	} else {
		church.EstablishedDate = nil
	}

	if err := uc.churchRepo.Update(church); err != nil {
		return nil, fmt.Errorf("failed to update church: %w", err)
	}

	return uc.entityToResponse(church), nil
}

func (uc *churchUseCase) DeleteChurch(id uint, userID uint) error {
	church, err := uc.churchRepo.FindByID(id)
	if err != nil {
		return errors.New("church not found")
	}

	// Check if user is the owner
	if church.OwnerID == nil || *church.OwnerID != userID {
		return errors.New("only church owner can delete the church")
	}

	return uc.churchRepo.Delete(id)
}

func (uc *churchUseCase) JoinChurch(userID uint, churchID uint) error {
	// Check if church exists
	_, err := uc.churchRepo.FindByID(churchID)
	if err != nil {
		return errors.New("church not found")
	}

	// Get user
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Check if user is already in a church
	if user.ChurchID != nil {
		return errors.New("user is already a member of a church. Please leave current church first")
	}

	// Add user to church with pending status
	user.ChurchID = &churchID
	user.ChurchStatus = entity.ChurchStatusPending

	return uc.userRepo.Update(user)
}

func (uc *churchUseCase) LeaveChurch(userID uint) error {
	user, err := uc.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	if user.ChurchID == nil {
		return errors.New("user is not a member of any church")
	}

	// Remove from church
	user.ChurchID = nil
	user.ChurchStatus = entity.ChurchStatusPending // Reset status

	return uc.userRepo.Update(user)
}

func (uc *churchUseCase) ApproveMember(churchID uint, ownerID uint, targetUserID uint, status string) error {
	// Check if church exists
	church, err := uc.churchRepo.FindByID(churchID)
	if err != nil {
		return errors.New("church not found")
	}

	// Verify owner
	if church.OwnerID == nil || *church.OwnerID != ownerID {
		return errors.New("only church owner can approve members")
	}

	// Get target user
	user, err := uc.userRepo.FindByID(targetUserID)
	if err != nil {
		return errors.New("user not found")
	}

	// Verify user is requesting to join this church
	if user.ChurchID == nil || *user.ChurchID != churchID {
		return errors.New("user is not requesting to join this church")
	}

	// Update status
	if status == "approved" {
		user.ChurchStatus = entity.ChurchStatusApproved
	} else if status == "rejected" {
		user.ChurchStatus = entity.ChurchStatusRejected
		user.ChurchID = nil // Remove from church if rejected
	} else {
		return errors.New("invalid status. Must be 'approved' or 'rejected'")
	}

	return uc.userRepo.Update(user)
}

func (uc *churchUseCase) GetPendingMembers(churchID uint, ownerID uint, page, pageSize int) (*dto.PaginatedResponse, error) {
	// Check if church exists
	church, err := uc.churchRepo.FindByID(churchID)
	if err != nil {
		return nil, errors.New("church not found")
	}

	// Verify owner
	if church.OwnerID == nil || *church.OwnerID != ownerID {
		return nil, errors.New("only church owner can view pending members")
	}

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
	users, err := uc.churchRepo.FindMembers(churchID, "pending", pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch pending members: %w", err)
	}

	responses := make([]interface{}, len(users))
	for i, user := range users {
		responses[i] = uc.userToResponse(user)
	}

	totalCount, err := uc.churchRepo.CountMembers(churchID, "pending")
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

func (uc *churchUseCase) GetApprovedMembers(churchID uint, page, pageSize int) (*dto.PaginatedResponse, error) {
	// Check if church exists
	_, err := uc.churchRepo.FindByID(churchID)
	if err != nil {
		return nil, errors.New("church not found")
	}

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
	users, err := uc.churchRepo.FindMembers(churchID, "approved", pageSize, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch approved members: %w", err)
	}

	responses := make([]interface{}, len(users))
	for i, user := range users {
		responses[i] = uc.userToResponse(user)
	}

	totalCount, err := uc.churchRepo.CountMembers(churchID, "approved")
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

func (uc *churchUseCase) entityToResponse(church *entity.Church) *dto.ChurchResponse {
	response := &dto.ChurchResponse{
		ID:              church.ID,
		Name:        church.Name,
		Address:         church.Address,
		Phone:           church.Phone,
		Email:           church.Email,
		Website:         church.Website,
		PastorName:      church.PastorName,
		Description:     church.Description,
		Logo:            church.Logo,
		EstablishedDate: dto.NewLocalTimePtr(church.EstablishedDate),
		Denomination:    church.Denomination,
		OwnerID:         church.OwnerID,
		CreatedAt:       dto.NewLocalTime(church.CreatedAt),
		UpdatedAt:       dto.NewLocalTime(church.UpdatedAt),
	}

	// Get member counts
	if approvedCount, err := uc.churchRepo.CountMembers(church.ID, "approved"); err == nil {
		response.MemberCount = int(approvedCount)
	}
	if pendingCount, err := uc.churchRepo.CountMembers(church.ID, "pending"); err == nil {
		response.PendingCount = int(pendingCount)
	}

	// Convert owner if present
	if church.Owner != nil {
		response.Owner = &dto.UserResponse{
			ID:       church.Owner.ID,
			Username: church.Owner.Username,
			Name: church.Owner.Name,
			Email:    church.Owner.Email,
			Profile:  church.Owner.Profile,
		}
	}

	return response
}

func (uc *churchUseCase) userToResponse(user *entity.User) *dto.UserResponse {
	return &dto.UserResponse{
		ID:           user.ID,
		Username:     user.Username,
		Name:     user.Name,
		Email:        user.Email,
		Profile:      user.Profile,
		ChurchID:     user.ChurchID,
		ChurchStatus: string(user.ChurchStatus),
		Birthday:     dto.NewLocalTimePtr(user.Birthday),
		Bio:          user.Bio,
		CreatedAt:    dto.NewLocalTime(user.CreatedAt),
		UpdatedAt:    dto.NewLocalTime(user.UpdatedAt),
	}
}

