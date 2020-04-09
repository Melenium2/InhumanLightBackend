package notifications

type NotificationSender interface{
	Notify() chan string
}