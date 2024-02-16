package notifications

import "github.com/mblancoa/authentication/core/ports"

type notificationService struct {
}

func NewNotificationService() ports.NotificationService {
	return &notificationService{}
}

func (n *notificationService) SendSMS(template string, source any, numberPhone string) {
	//TODO implement me
	panic("implement me")
}

func (n *notificationService) SendEmail(template string, source any, email string) {
	//TODO implement me
	panic("implement me")
}
