package events

import (
	"prc_hub-api/users"
	"time"
)

type Event struct {
	Id                  uint64            `json:"id"`
	UserId              uint64            `json:"user_id,omitempty"`
	Title               string            `json:"title"`
	Description         *string           `json:"description,omitempty"`
	Speakers            []users.UserEmbed `json:"speakers,omitempty"`
	Location            string            `json:"location,omitempty"`
	Datetimes           []EventDatetime   `json:"datetimes"`
	Published           bool              `json:"published"`
	Completed           bool              `json:"completed"`
	AutoNotifyDocuments bool              `json:"auto_notify_documents_enabled"`
	Documents           []EventDocument   `json:"documents,omitempty"`
}

type EventDocument struct {
	Id   uint64 `json:"id"`
	Name string `json:"name"`
	Url  string `json:"url"`
}

type EventDatetime struct {
	Id      uint64     `json:"id"`
	EventId uint64     `json:"event_id"`
	Start   time.Time  `json:"start"`
	End     *time.Time `json:"end"`
}

type EventParticipate struct {
	EventDatetimeId uint64 `json:"event_datetime_id"`
	UserId          uint64 `json:"user_id"`
}
