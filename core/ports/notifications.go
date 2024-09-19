package ports

var NotificationContext *notificationContext = &notificationContext{}

type notificationContext struct {
	NotificationService NotificationService
}

type NotificationService interface {
	SendSMS(template string, source any, numberPhone string)
	SendEmail(template string, source any, email string)
}
