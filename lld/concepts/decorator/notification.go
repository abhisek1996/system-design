package decorator

import "fmt"

type Notification interface {
	Notify(message string)
	Cost() int
}

type EmailNotification struct {
	emailId string
}

func (e *EmailNotification) Notify(message string) {
	fmt.Println("Email Notification: " + message)
}

func (e *EmailNotification) Cost() int {
	return 1
}

// we want to send via sms, push, mobile
// permutation and combination og the above

// NotificationDecorator is a decorator for Notification

type NotificationDecorator struct { // we can directly embed Notification without this struct
	Notification
}

type MobileNotification struct {
	NotificationDecorator
	mobileNumber string
}

func NewMobileNotification(notification Notification) Notification {
	return &MobileNotification{NotificationDecorator: NotificationDecorator{Notification: notification}}
}

func (m *MobileNotification) Notify(message string) {
	m.NotificationDecorator.Notify(message)
	fmt.Println("Mobile Notification: " + message)
}

type SMSNotification struct {
	NotificationDecorator
	mobileNumber string
}

func NewSMSNotification(notification Notification) Notification {
	return &SMSNotification{NotificationDecorator: NotificationDecorator{Notification: notification}}
}

func (s *SMSNotification) Notify(message string) {
	s.NotificationDecorator.Notify(message)
	fmt.Println("SMS Notification: " + message)
}

func NotificationDemo() {
	emailNotification := &EmailNotification{emailId: "abc@abc.com"}
	//	emailNotification.Notify("Hello")

	mobileNotification := NewMobileNotification(emailNotification)
	mobileNotification.Notify("Hello")

	fmt.Println("------------")

	smsNotification := NewSMSNotification(mobileNotification)
	smsNotification.Notify("Hello")
}
