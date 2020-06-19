// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package operation

import "github.com/brannon/apnstool/apns"

type SendRawOperation struct {
	SendOperation

	DataString string
	Priority   string
	PushType   string
}

func NewSendRaw(client apns.Client) *SendRawOperation {
	op := &SendRawOperation{}
	op.Client = client
	return op
}

func (op *SendRawOperation) Exec() (*SendOperationResult, error) {
	headers := make(apns.Headers)

	headers["apns-topic"] = op.AppId

	if op.Priority != "" {
		headers["apns-priority"] = op.Priority
	}

	if op.PushType != "" {
		headers["apns-push-type"] = op.PushType
	}

	return op.sendNotification(headers, []byte(op.DataString))
}
