package services

import (
	"context"

	"gorm.io/gorm"

	"github.com/reshap0318/serv-uam/internal/dtos"
	"github.com/reshap0318/serv-uam/internal/helpers"
	"github.com/reshap0318/serv-uam/internal/models"
)

// ProfileGetMe returns the caller's own profile (identity from GatewayAuth, see
// helpers.GetCallerID).
func (s *Services) ProfileGetMe(ctx context.Context) (*dtos.ProfileDTO, error) {
	user, err := s.repo.User.FindByID(nil, helpers.GetCallerID(ctx), "Roles.Permissions")
	if err != nil {
		return nil, err
	}
	dto := dtos.ToProfileDTO(user)
	return &dto, nil
}

// ProfileUpdateMe updates the caller's own name/email/password. Unlike UserUpdate, roles and
// status are never touched here — a user can't grant themselves roles or reactivate a
// suspended account this way.
func (s *Services) ProfileUpdateMe(ctx context.Context, req dtos.ProfileUpdateRequest) (*dtos.ProfileDTO, error) {
	id := helpers.GetCallerID(ctx)

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
		Email: req.Email,
		Name:  req.Name,
	}
	if req.Password != "" {
		hashed, err := helpers.HashString(req.Password)
		if err != nil {
			s.Logger.LogCtx(ctx, "ProfileUpdateMe", "Failed to hash password: %v", err)
			return nil, err
		}
		update.Password = hashed
	}

	var result *models.User
	err = s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		var err error
		result, err = s.repo.User.Update(tx, &models.User{ID: id}, update)
		return err
	})
	if err != nil {
		s.Logger.LogCtx(ctx, "ProfileUpdateMe", "Failed: %v", err)
		return nil, err
	}

	result, err = s.repo.User.FindByID(nil, id, "Roles.Permissions")
	if err != nil {
		return nil, err
	}

	s.Logger.LogCtx(ctx, "ProfileUpdateMe", "User %d updated own profile", id)
	dto := dtos.ToProfileDTO(result)
	return &dto, nil
}
