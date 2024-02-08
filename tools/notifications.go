package tools

type NotificationService interface {
	SendSMS(template string, source any, numberPhone string)
	SendEmail(template string, source any, email string)
}

type notificationService struct {
}

func NewNotificationService() NotificationService {
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
