// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package operation

import (
	"github.com/brannon/apnstool/apns"
)

type SendAlertOperation struct {
	SendOperation

	Expiration string
	Priority   string

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

	if op.Expiration != "" {
		expiration, err := apns.ParseExpiration(op.Expiration)
		if err != nil {
			return nil, err
		}
		notificationBuilder.SetExpiration(expiration)
	}

	if op.Priority != "" {
		priority, err := apns.ParsePriority(op.Priority)
		if err != nil {
			return nil, err
		}
		notificationBuilder.SetPriority(priority)
	}

	headers, content, err := notificationBuilder.Build()
	if err != nil {
		return nil, err
	}

	return op.sendNotification(headers, content)
}
