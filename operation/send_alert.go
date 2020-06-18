package operation

import "github.com/brannon/apnstool/apns"

type SendAlertOperation struct {
	SendOperation

	AlertText  string
	BadgeCount int
	SoundName  string
}

func NewSendAlert(client apns.Client) *SendAlertOperation {
	op := &SendAlertOperation{}
	op.Client = client
	return op
}

func (op *SendAlertOperation) Exec() (*SendOperationResult, error) {
	notificationBuilder := apns.NewNotificationBuilder(op.AppId)

	if op.AlertText != "" {
		notificationBuilder.SetAlertText(op.AlertText)
	}

	if op.BadgeCount != -1 {
		notificationBuilder.SetBadgeCount(op.BadgeCount)
	}

	if op.SoundName != "" {
		notificationBuilder.SetSoundName(op.SoundName)
	}

	headers, content, err := notificationBuilder.Build()
	if err != nil {
		return nil, err
	}

	return op.sendNotification(headers, content)
}
