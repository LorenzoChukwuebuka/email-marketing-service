package payloads

type AdminNotificationPayload struct {
	UserId            string
	Link              string
	NotificationTitle string
}

type UserNotificationPayload struct {
	UserId           string
	NotifcationTitle string
}

