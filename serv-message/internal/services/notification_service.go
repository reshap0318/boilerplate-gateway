package services

import (
	"context"
	"time"

	"gorm.io/gorm"

	"github.com/reshap0318/serv-message/internal/dtos"
	"github.com/reshap0318/serv-message/internal/helpers"
	"github.com/reshap0318/serv-message/internal/models"
	"github.com/reshap0318/serv-message/internal/repositories"
)

// NotificationCreate creates a notification for req.UserID. senderID is 0
// when the call is system-triggered — this endpoint has no auth middleware
// (see routes/notification_route.go), so the handler passes whatever
// X-User-Id it found on the request, if any.
func (s *Services) NotificationCreate(ctx context.Context, req dtos.NotificationCreateRequest, senderID uint) (*dtos.NotificationDTO, error) {
	s.Logger.LogCtx(ctx, "NotificationCreate", "Creating notification for user %d: %s", req.UserID, req.Title)

	notification := &models.Notification{
		UserID:   req.UserID,
		SenderID: senderID,
		Type:     req.Type,
		Title:    req.Title,
		Message:  req.Message,
		Data:     req.Data,
	}

	res, err := s.repo.TxManager.WithinTransactionWithResult(func(tx *gorm.DB) (interface{}, error) {
		return s.repo.Notification.Create(tx, notification)
	})
	if err != nil {
		s.Logger.LogCtx(ctx, "NotificationCreate", "Failed to create notification: %v", err)
		return nil, err
	}

	dto := dtos.ToNotificationDTO(res.(*models.Notification))
	return &dto, nil
}

// NotificationGetAll returns paginated notifications belonging to the caller.
func (s *Services) NotificationGetAll(ctx context.Context, opts *repositories.QueryOptions) (*repositories.PagedResult[dtos.NotificationDTO], error) {
	if opts == nil {
		opts = &repositories.QueryOptions{}
	}
	if opts.SortBy == "" {
		opts.SortBy = "id"
	}
	if opts.Order == "" {
		opts.Order = "DESC"
	}
	opts.ConditionGroups = append(opts.ConditionGroups, repositories.ConditionGroup{
		Conditions: []repositories.QueryCondition{
			{Column: "user_id", Operator: "=", Value: helpers.GetCallerID(ctx)},
		},
	})

	result, err := s.repo.Notification.FindAllWithOpts(nil, opts)
	if err != nil {
		return nil, err
	}

	return &repositories.PagedResult[dtos.NotificationDTO]{
		Data:       dtos.ToNotificationDTOList(result.Data),
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// NotificationCountUnread returns the number of unread notifications for the caller.
func (s *Services) NotificationCountUnread(ctx context.Context) (int64, error) {
	return s.repo.Notification.CountUnread(nil, helpers.GetCallerID(ctx))
}

// NotificationMarkAsRead marks a single notification (owned by the caller) as read.
func (s *Services) NotificationMarkAsRead(ctx context.Context, id uint) error {
	userID := helpers.GetCallerID(ctx)

	return s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		_, err := s.repo.Notification.UpdateMap(tx, &models.Notification{ID: id, UserID: userID}, map[string]interface{}{"read_at": time.Now()})
		return err
	})
}

// NotificationMarkAllAsRead marks every unread notification of the caller as read.
func (s *Services) NotificationMarkAllAsRead(ctx context.Context) error {
	userID := helpers.GetCallerID(ctx)

	return s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		return s.repo.Notification.MarkAllAsRead(tx, userID)
	})
}

// NotificationDelete soft-deletes a single notification owned by the caller.
func (s *Services) NotificationDelete(ctx context.Context, id uint) error {
	userID := helpers.GetCallerID(ctx)

	return s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		return s.repo.Notification.DeleteOwned(tx, id, userID)
	})
}

// NotificationDeleteAll soft-deletes every notification belonging to the caller.
func (s *Services) NotificationDeleteAll(ctx context.Context) error {
	userID := helpers.GetCallerID(ctx)

	return s.repo.TxManager.WithinTransaction(func(tx *gorm.DB) error {
		return s.repo.Notification.DeleteAllByUser(tx, userID)
	})
}
