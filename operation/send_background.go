// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package operation

import "github.com/brannon/apnstool/apns"

type SendBackgroundOperation struct {
	SendOperation

	DataString string
}

func NewSendBackground(client apns.Client) *SendBackgroundOperation {
	op := &SendBackgroundOperation{}
	op.Client = client
	return op
}

func (op *SendBackgroundOperation) Exec() (*SendOperationResult, error) {
	notificationBuilder := apns.NewNotificationBuilder(op.AppId)
	notificationBuilder.SetContentAvailable(true)

	if op.DataString != "" {
		data, err := parseDataString(op.DataString)
		if err != nil {
			return nil, err
		}

		notificationBuilder.Merge(data)
	}

	headers, content, err := notificationBuilder.Build()
	if err != nil {
		return nil, err
	}

	return op.sendNotification(headers, content)
}
