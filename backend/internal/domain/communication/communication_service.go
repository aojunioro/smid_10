package communication

import "context"

type NotificationService struct {
	notificationRepo NotificationRepository
}

func NewNotificationService(notificationRepo NotificationRepository) *NotificationService {
	return &NotificationService{
		notificationRepo: notificationRepo,
	}
}

type CreateNotificationRequest struct {
	SystemUserID   *int64  `json:"system_user_id"`
	SystemUserToID *int64  `json:"system_user_to_id"`
	Subject        string  `json:"subject"`
	Message        string  `json:"message"`
	ActionURL      *string `json:"action_url"`
	ActionLabel    *string `json:"action_label"`
	Icon           *string `json:"icon"`
}

type UpdateNotificationRequest struct {
	Subject     *string `json:"subject"`
	Message     *string `json:"message"`
	ActionURL   *string `json:"action_url"`
	ActionLabel *string `json:"action_label"`
	Icon        *string `json:"icon"`
	Checked     *bool   `json:"checked"`
}

func (s *NotificationService) Create(ctx context.Context, req CreateNotificationRequest) (*SystemNotification, error) {
	notification := &SystemNotification{
		SystemUserID:   req.SystemUserID,
		SystemUserToID: req.SystemUserToID,
		Subject:        req.Subject,
		Message:        req.Message,
		ActionURL:      req.ActionURL,
		ActionLabel:    req.ActionLabel,
		Icon:           req.Icon,
	}

	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *NotificationService) GetByID(ctx context.Context, id int64) (*SystemNotification, error) {
	notification, err := s.notificationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return notification, nil
}

func (s *NotificationService) List(ctx context.Context, limit, offset int) ([]SystemNotification, error) {
	opts := ListOptions{
		Limit:  limit,
		Offset: offset,
	}

	notifications, err := s.notificationRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return notifications, nil
}

func (s *NotificationService) Update(ctx context.Context, id int64, req UpdateNotificationRequest) (*SystemNotification, error) {
	notification, err := s.notificationRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Subject != nil {
		notification.Subject = *req.Subject
	}
	if req.Message != nil {
		notification.Message = *req.Message
	}
	if req.ActionURL != nil {
		notification.ActionURL = req.ActionURL
	}
	if req.ActionLabel != nil {
		notification.ActionLabel = req.ActionLabel
	}
	if req.Icon != nil {
		notification.Icon = req.Icon
	}
	if req.Checked != nil {
		notification.Checked = req.Checked
	}

	if err := s.notificationRepo.Update(ctx, notification); err != nil {
		return nil, err
	}

	return notification, nil
}

type MessageService struct {
	messageRepo MessageRepository
}

func NewMessageService(messageRepo MessageRepository) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
	}
}

type CreateMessageRequest struct {
	SystemUserID   *int64  `json:"system_user_id"`
	SystemUserToID *int64  `json:"system_user_to_id"`
	Subject        string  `json:"subject"`
	Message        string  `json:"message"`
	Attachments    *string `json:"attachments"`
}

type UpdateMessageRequest struct {
	Subject     *string `json:"subject"`
	Message     *string `json:"message"`
	Checked     *bool   `json:"checked"`
	Removed     *bool   `json:"removed"`
	Viewed      *bool   `json:"viewed"`
	Attachments *string `json:"attachments"`
}

func (s *MessageService) Create(ctx context.Context, req CreateMessageRequest) (*SystemMessage, error) {
	message := &SystemMessage{
		SystemUserID:   req.SystemUserID,
		SystemUserToID: req.SystemUserToID,
		Subject:        req.Subject,
		Message:        req.Message,
		Attachments:    req.Attachments,
	}

	if err := s.messageRepo.Create(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}

func (s *MessageService) GetByID(ctx context.Context, id int64) (*SystemMessage, error) {
	message, err := s.messageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return message, nil
}

func (s *MessageService) List(ctx context.Context, limit, offset int) ([]SystemMessage, error) {
	opts := ListOptions{
		Limit:  limit,
		Offset: offset,
	}

	messages, err := s.messageRepo.List(ctx, opts)
	if err != nil {
		return nil, err
	}

	return messages, nil
}

func (s *MessageService) Update(ctx context.Context, id int64, req UpdateMessageRequest) (*SystemMessage, error) {
	message, err := s.messageRepo.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if req.Subject != nil {
		message.Subject = *req.Subject
	}
	if req.Message != nil {
		message.Message = *req.Message
	}
	if req.Checked != nil {
		message.Checked = req.Checked
	}
	if req.Removed != nil {
		message.Removed = req.Removed
	}
	if req.Viewed != nil {
		message.Viewed = req.Viewed
	}
	if req.Attachments != nil {
		message.Attachments = req.Attachments
	}

	if err := s.messageRepo.Update(ctx, message); err != nil {
		return nil, err
	}

	return message, nil
}
