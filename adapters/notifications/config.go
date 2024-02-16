package notifications

import "github.com/mblancoa/authentication/core"

func SetupNotificationConfiguration() {
	setupNotificationContext()
}
func setupNotificationContext() {
	ctx := core.NotificationContext
	ctx.NotificationService = NewNotificationService()
}
