// Copyright 2019 Brannon Jones. All rights reserved.
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package apns

import "encoding/json"

type NotificationBuilder struct {
	AppId   string
	content map[string]interface{}
}

func NewNotificationBuilder(appId string) *NotificationBuilder {
	return &NotificationBuilder{
		AppId:   appId,
		content: make(map[string]interface{}),
	}
}

func (b *NotificationBuilder) aps() map[string]interface{} {
	aps, ok := b.content["aps"].(map[string]interface{})
	if !ok {
		aps = make(map[string]interface{})
		b.content["aps"] = aps
	}
	return aps
}

func (b *NotificationBuilder) Build() (Headers, []byte, error) {
	headers, err := b.BuildHeaders()
	if err != nil {
		return nil, nil, err
	}

	content, err := b.BuildContent()
	if err != nil {
		return nil, nil, err
	}

	return headers, content, nil
}

func (b *NotificationBuilder) BuildContent() ([]byte, error) {
	data, err := json.Marshal(b.content)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (b *NotificationBuilder) BuildHeaders() (Headers, error) {
	headers := Headers{}

	headers["apns-topic"] = b.AppId

	pushType := b.getPushType()
	headers["apns-push-type"] = pushType

	if pushType == "background" {
		headers["apns-priority"] = "5"
	}

	return headers, nil
}

func (b *NotificationBuilder) getPushType() string {
	aps := b.aps()
	if hasKey(aps, "alert") || hasKey(aps, "badge") || hasKey(aps, "sound") {
		return "alert"
	} else if hasKey(aps, "content-available") && aps["content-available"] == 1 {
		return "background"
	}
	return ""
}

func (b *NotificationBuilder) Merge(data map[string]interface{}) *NotificationBuilder {
	for k, v := range data {
		b.content[k] = v
	}
	return b
}

func (b *NotificationBuilder) SetAlertText(text string) *NotificationBuilder {
	b.aps()["alert"] = text
	return b
}

func (b *NotificationBuilder) SetBadgeCount(count int) *NotificationBuilder {
	b.aps()["badge"] = count
	return b
}

func (b *NotificationBuilder) SetContentAvailable(value bool) *NotificationBuilder {
	var intValue int = 0
	if value {
		intValue = 1
	}

	b.aps()["content-available"] = intValue
	return b
}

func (b *NotificationBuilder) SetSoundName(name string) *NotificationBuilder {
	b.aps()["sound"] = name
	return b
}

func hasKey(m map[string]interface{}, name string) bool {
	_, ok := m[name]
	return ok
}
