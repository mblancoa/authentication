package notifications

import (
	"github.com/mblancoa/authentication/core/ports"
)

func SetupNotificationConfiguration() {
	setupNotificationContext()
}
func setupNotificationContext() {
	ctx := ports.NotificationContext
	ctx.NotificationService = NewNotificationService()
}
