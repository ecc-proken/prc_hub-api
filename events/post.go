package events

import (
	"prc_hub-api/mysql"
	"prc_hub-api/users"
	"strings"
	"time"
)

type PostBody struct {
	Title               string         `json:"title" validate:"required,gte=1"`
	Description         *string        `json:"description" validate:"omitempty"`
	Speakers            []uint64       `json:"speakers" validate:"required,gte=1,dive,gte=1"`
	Location            string         `json:"location" validate:"omitempty,gte=1"`
	Datetimes           []PostDatetime `json:"datetimes" validate:"required,gte=1,dive"`
	Published           bool           `json:"published"`
	Completed           bool           `json:"completed"`
	AutoNotifyDocuments bool           `json:"auto_notify_documents_enabled"`
}

type PostDatetime struct {
	Start time.Time  `json:"start" validate:"required"`
	End   *time.Time `json:"end" validate:"omitempty,gtcsfield=Start"`
}

func Post(userId uint64, post PostBody) (e Event, notFoundUserIds []uint64, err error) {
	// DBに接続
	db, err := mysql.Open()
	if err != nil {
		return
	}
	// トランザクション開始
	defer db.Close()
	tx, err := db.Begin()
	if err != nil {
		return
	}

	// 書込
	result1, err := mysql.TxWrite(
		tx,
		`INSERT INTO events (user_id, title, description, location, published, completed, auto_notify_documents_enabled)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
		userId, post.Title, post.Description, post.Location, post.Published, post.Completed, post.AutoNotifyDocuments,
	)
	if err != nil {
		tx.Rollback()
		return
	}
	// Insertした行のIdを取得
	eventId, err := result1.LastInsertId()
	if err != nil {
		tx.Rollback()
		return
	}

	// クエリを作成
	queryStr2 := "INSERT INTO event_speakers (event_id, user_id) VALUES "
	var queryParams2 []interface{}
	for _, v := range post.Speakers {
		queryStr2 += "(?, ?),"
		queryParams2 = append(queryParams2, eventId, v)
	}
	queryStr2 = strings.TrimSuffix(queryStr2, ",")
	// 書込
	_, err = mysql.TxWrite(
		tx,
		queryStr2,
		queryParams2...,
	)
	if err != nil {
		tx.Rollback()
		return
	}
	// ユーザー取得と確認
	eventSpeakers, err := users.GetEmbed(users.GetEmbedQuery{Ids: post.Speakers})
	if err != nil {
		tx.Rollback()
		return
	}
	if len(post.Speakers) != len(eventSpeakers) {
		// 不正なIdが指定された場合、Idを特定
		for _, id := range post.Speakers {
			found := false
			for _, user := range eventSpeakers {
				if id == user.Id {
					found = true
					break
				}
			}
			if !found {
				notFoundUserIds = append(notFoundUserIds, id)
			}
		}
	}

	// クエリを作成
	queryStr3 := "INSERT INTO event_datetimes (event_id, start, end) VALUES "
	var queryParams3 []interface{}
	for _, v := range post.Datetimes {
		queryStr3 += "(?, ?, ?),"
		queryParams3 = append(queryParams3, eventId, v.Start, v.End)
	}
	queryStr3 = strings.TrimSuffix(queryStr3, ",")
	// 書込
	result3, err := mysql.TxWrite(
		tx,
		queryStr3,
		queryParams3...,
	)
	if err != nil {
		tx.Rollback()
		return
	}
	eventDatetimeId, err := result3.LastInsertId()
	if err != nil {
		tx.Rollback()
		return
	}
	var eventDatetimes []EventDatetime
	for i, d := range post.Datetimes {
		eventDatetimes = append(
			eventDatetimes,
			EventDatetime{
				Id:      uint64(eventDatetimeId + int64(i)),
				EventId: uint64(eventId),
				Start:   d.Start,
				End:     d.End,
			},
		)
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	e.Id = uint64(eventId)
	e.UserId = userId
	e.Title = post.Title
	e.Description = post.Description
	e.Location = post.Location
	e.Published = post.Published
	e.Completed = post.Completed
	e.AutoNotifyDocuments = post.AutoNotifyDocuments
	e.Speakers = eventSpeakers
	e.Datetimes = eventDatetimes

	return
}
