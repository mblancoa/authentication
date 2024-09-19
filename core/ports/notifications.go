package ports

type NotificationService interface {
	SendSMS(template string, source any, numberPhone string)
	SendEmail(template string, source any, email string)
}
