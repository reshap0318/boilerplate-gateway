package services

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/reshap0318/serv-uam/internal/dtos"
	"github.com/reshap0318/serv-uam/internal/helpers"
	"github.com/reshap0318/serv-uam/internal/models"
	"github.com/reshap0318/serv-uam/internal/repositories"
)

// UserCreate creates a new user and attaches the given roles.
func (s *Services) UserCreate(ctx context.Context, req dtos.UserCreateRequest) (*dtos.UserDTO, error) {
	s.Logger.LogCtx(ctx, "UserCreate", "Creating user: %s", req.Email)

	exists, err := s.repo.User.ExistsWithMap(nil, map[string]interface{}{"email": req.Email})
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, &helpers.FieldError{Field: "email", Message: "user already exists"}
	}

	hashed, err := helpers.HashString(req.Password)
	if err != nil {
		s.Logger.LogCtx(ctx, "UserCreate", "Failed to hash password: %v", err)
		return nil, err
	}

	status := req.Status
	if status == "" {
		status = "active"
	}

	var result *models.User
	err = s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		user := &models.User{
			Email:    req.Email,
			Password: hashed,
			Name:     req.Name,
			Status:   status,
		}
		var err error
		result, err = s.repo.User.Create(tx, user)
		if err != nil {
			return err
		}

		return s.repo.User.SyncRoles(tx, result, req.Roles)
	})
	if err != nil {
		s.Logger.LogCtx(ctx, "UserCreate", "Failed: %v", err)
		return nil, err
	}

	result, err = s.repo.User.FindByID(nil, result.ID, "Roles.Permissions")
	if err != nil {
		return nil, err
	}

	s.Logger.LogCtx(ctx, "UserCreate", "User created: %s (ID: %d)", result.Email, result.ID)
	dto := dtos.ToUserDTO(result)
	return &dto, nil
}

// UserGetAll returns a paginated list of users with their roles.
func (s *Services) UserGetAll(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.UserDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	opts.Preloads = append(opts.Preloads, "Roles.Permissions")

	result, err := s.repo.User.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	return &repositories.PagedResult[dtos.UserDTO]{
		Data:       dtos.ToUserDTOList(result.Data),
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// UserGetByID returns a single user with their roles.
func (s *Services) UserGetByID(ctx context.Context, id uint) (*dtos.UserDTO, error) {
	user, err := s.repo.User.FindByID(nil, id, "Roles.Permissions")
	if err != nil {
		return nil, err
	}
	dto := dtos.ToUserDTO(user)
	return &dto, nil
}

// UserUpdate updates a user and replaces their role set. Password is left
// unchanged when req.Password is empty.
func (s *Services) UserUpdate(ctx context.Context, id uint, req dtos.UserUpdateRequest) (*dtos.UserDTO, error) {
	s.Logger.LogCtx(ctx, "UserUpdate", "Updating user %d", id)

	existing, err := s.repo.User.FindByID(nil, id)
	if err != nil {
		return nil, err
	}

	if existing.Email != req.Email {
		exists, err := s.repo.User.ExistsWithMap(nil, map[string]interface{}{"email": req.Email})
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, &helpers.FieldError{Field: "email", Message: "user already exists"}
		}
	}

	update := &models.User{
		Email:  req.Email,
		Name:   req.Name,
		Status: req.Status,
	}
	if req.Password != "" {
		hashed, err := helpers.HashString(req.Password)
		if err != nil {
			s.Logger.LogCtx(ctx, "UserUpdate", "Failed to hash password: %v", err)
			return nil, err
		}
		update.Password = hashed
	}

	var result *models.User
	err = s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.repo.User.Update(tx, &models.User{ID: id}, update)
		if err != nil {
			return err
		}

		return s.repo.User.SyncRoles(tx, result, req.Roles)
	})
	if err != nil {
		s.Logger.LogCtx(ctx, "UserUpdate", "Failed: %v", err)
		return nil, err
	}

	result, err = s.repo.User.FindByID(nil, id, "Roles.Permissions")
	if err != nil {
		return nil, err
	}

	s.invalidateUserAccess(id)

	s.Logger.LogCtx(ctx, "UserUpdate", "User %d updated", id)
	dto := dtos.ToUserDTO(result)
	return &dto, nil
}

// UserUnlock manually clears a user's account lock before it would
// naturally expire (see registerFailedLogin in auth_service.go).
func (s *Services) UserUnlock(ctx context.Context, id uint) (*dtos.UserDTO, error) {
	s.Logger.LogCtx(ctx, "UserUnlock", "Unlocking user %d", id)

	existing, err := s.repo.User.FindByID(nil, id)
	if err != nil {
		return nil, err
	}

	if existing.LockedUntil == nil || !existing.LockedUntil.After(time.Now()) {
		return nil, &helpers.CustomError{Status: 400, Message: "This account is not currently locked"}
	}

	if _, err := s.repo.User.UpdateMap(nil, &models.User{ID: id}, map[string]interface{}{
		"locked_until":          nil,
		"failed_login_attempts": 0,
	}); err != nil {
		s.Logger.LogCtx(ctx, "UserUnlock", "Failed: %v", err)
		return nil, err
	}

	result, err := s.repo.User.FindByID(nil, id, "Roles.Permissions")
	if err != nil {
		return nil, err
	}

	s.Logger.LogCtx(ctx, "UserUnlock", "User %d unlocked", id)
	dto := dtos.ToUserDTO(result)
	return &dto, nil
}

// UserDelete deletes a user.
func (s *Services) UserDelete(ctx context.Context, id uint) error {
	s.Logger.LogCtx(ctx, "UserDelete", "Deleting user %d", id)

	err := s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		_, err := s.repo.User.Delete(tx, id)
		return err
	})
	if err != nil {
		s.Logger.LogCtx(ctx, "UserDelete", "Failed: %v", err)
		return err
	}

	s.invalidateUserAccess(id)

	s.Logger.LogCtx(ctx, "UserDelete", "User %d deleted", id)
	return nil
}
