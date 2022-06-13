package events

import (
	"fmt"
	"prc_hub-api/webhook"
	"strings"
)

func (e *Event) NotifyLINE(token string, frontEndUrl string) (err error) {
	// メッセージを生成
	msg := "イベント情報\n\n%s: %s\n\n"
	msgParams := []interface{}{"勉強会", e.Title}
	if e.Description != nil {
		msg += "%s\n\n"
		msgParams = append(msgParams, *e.Description)
	}
	msg += "%s/events/%d"
	msgParams = append(msgParams, frontEndUrl, e.Id)

	webhook.LineNotify(token, fmt.Sprintf(msg, msgParams...))
	return
}

func (e *Event) NotifyLINEDocuments(token string) (err error) {
	if e.Documents != nil && len(e.Documents) != 0 {
		// メッセージを生成
		msg := "イベント資料\n\n%s: %s\n\n"
		msgParams := []interface{}{"勉強会", e.Title}

		for _, d := range e.Documents {
			msg += "%s\n%s\n\n"
			msgParams = append(msgParams, d.Name, d.Url)
		}
		msg = strings.TrimSuffix(msg, "\n\n")

		webhook.LineNotify(token, fmt.Sprintf(msg, msgParams...))
	}
	return
}
