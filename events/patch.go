package events

import (
	"database/sql"
	"prc_hub-api/mysql"
	"prc_hub-api/users"
	"strings"
)

type PatchBody struct {
	Title               *string                                         `json:"title" validate:"omitempty,gte=1"`
	Description         mysql.PatchNullJSONString                       `json:"description" validate:"omitempty,dive"`
	Speakers            *[]uint64                                       `json:"speakers" validate:"omitempty,gte=1,dive,gte=1"`
	Location            mysql.PatchNullJSONString                       `json:"location" validate:"omitempty,dive"`
	Datetimes           *[]PostDatetime                                 `json:"datetimes" validate:"omitempty,gte=1,dive"`
	Published           *bool                                           `json:"published" validate:"omitempty"`
	Completed           *bool                                           `json:"completed" validate:"omitempty"`
	AutoNotifyDocuments *bool                                           `json:"auto_notify_documents_enabled" validate:"omitempty"`
	Documents           mysql.PatchNullJSONSlice[PostEventDocumentBody] `json:"documents" validate:"omitempty,dive"`
}

func Patch(id uint64, p PatchBody) (e Event, notFound bool, notFoundUserIds []uint64, err error) {
	// Eventを取得
	e, notFound, err = GetById(id)
	if err != nil {
		return
	}
	if notFound {
		return
	}

	// トランザクション開始
	tx, err := mysql.Begin()
	if err != nil {
		return
	}

	updated := e

	// eventsテーブルを更新(PATCH)
	if p.Title != nil || p.Description.String != nil || p.Location.String != nil ||
		p.Published != nil || p.Completed != nil || p.AutoNotifyDocuments != nil {
		// クエリ作成
		queryStr1 := "UPDATE events SET"
		var queryParams1 []interface{}
		if p.Title != nil {
			queryStr1 += " title = ?,"
			queryParams1 = append(queryParams1, p.Title)
			updated.Title = *p.Title
		}
		if p.Description.String != nil {
			if *p.Description.String != nil {
				queryStr1 += " description = ?,"
				queryParams1 = append(queryParams1, **p.Description.String)
				updated.Description = *p.Description.String
			} else {
				queryStr1 += " description = ?,"
				queryParams1 = append(queryParams1, nil)
				updated.Description = nil
			}
		}
		if p.Location.String != nil {
			if *p.Location.String != nil {
				queryStr1 += " location = ?,"
				queryParams1 = append(queryParams1, **p.Location.String)
				updated.Location = *p.Location.String
			} else {
				queryStr1 += " location = ?,"
				queryParams1 = append(queryParams1, nil)
				updated.Location = nil
			}
		}
		if p.Published != nil {
			queryStr1 += " published = ?,"
			queryParams1 = append(queryParams1, p.Published)
			updated.Published = *p.Published
		}
		if p.Completed != nil {
			queryStr1 += " completed = ?,"
			queryParams1 = append(queryParams1, p.Completed)
			updated.Completed = *p.Completed
		}
		if p.AutoNotifyDocuments != nil {
			queryStr1 += " auto_notify_documents_enabled = ?,"
			queryParams1 = append(queryParams1, p.AutoNotifyDocuments)
			updated.AutoNotifyDocuments = *p.AutoNotifyDocuments
		}
		queryStr1 = strings.TrimRight(queryStr1, ",")
		queryStr1 += " WHERE id = ?"
		queryParams1 = append(queryParams1, id)

		// 更新
		_, err = mysql.TxWrite(tx, queryStr1, queryParams1...)
		if err != nil {
			return
		}
	}

	// event_speakersテーブルを更新(PUT)
	if p.Speakers != nil {
		// ユーザー取得と確認
		var eventSpeakers []users.UserEmbed
		eventSpeakers, err = users.GetEmbed(users.GetEmbedQuery{Ids: *p.Speakers})
		if err != nil {
			tx.Rollback()
			return
		}
		if len(*p.Speakers) != len(eventSpeakers) {
			// 不正なIdが指定された場合、Idを特定
			for _, id := range *p.Speakers {
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
			if len(notFoundUserIds) != 0 {
				tx.Rollback()
				return
			}
		}

		// 削除
		_, err = mysql.TxWrite(tx, "DELETE FROM event_speakers WHERE event_id = ?", id)
		if err != nil {
			tx.Rollback()
			return
		}

		// クエリを作成
		queryStr2 := "INSERT INTO event_speakers (event_id, user_id) VALUES "
		var queryParams2 []interface{}
		for _, v := range *p.Speakers {
			queryStr2 += "(?, ?),"
			queryParams2 = append(queryParams2, id, v)
		}
		queryStr2 = strings.TrimSuffix(queryStr2, ",")
		// 書込
		_, err = mysql.TxWrite(tx, queryStr2, queryParams2...)
		if err != nil {
			tx.Rollback()
			return
		}

		updated.Speakers = eventSpeakers
	}

	// event_datetimesテーブルを更新(PUT)
	if p.Datetimes != nil {
		// 削除
		_, err = mysql.TxWrite(tx, "DELETE FROM event_datetimes WHERE event_id = ?", id)
		if err != nil {
			tx.Rollback()
			return
		}

		// クエリを作成
		queryStr3 := "INSERT INTO event_datetimes (event_id, start, end) VALUES "
		var queryParams3 []interface{}
		for _, v := range *p.Datetimes {
			queryStr3 += "(?, ?, ?),"
			queryParams3 = append(queryParams3, id, v.Start, v.End)
		}
		queryStr3 = strings.TrimSuffix(queryStr3, ",")
		// 書込
		var result3 sql.Result
		result3, err = mysql.TxWrite(tx, queryStr3, queryParams3...)
		if err != nil {
			tx.Rollback()
			return
		}
		var eventDatetimeId int64
		eventDatetimeId, err = result3.LastInsertId()
		if err != nil {
			tx.Rollback()
			return
		}
		var eventDatetimes []EventDatetime
		for i, d := range *p.Datetimes {
			eventDatetimes = append(
				eventDatetimes,
				EventDatetime{
					Id:      uint64(eventDatetimeId + int64(i)),
					EventId: uint64(id),
					Start:   d.Start,
					End:     d.End,
				},
			)
		}

		updated.Datetimes = eventDatetimes
	}

	// event_datetimesテーブルを更新(PUT)
	if p.Documents.Slice != nil {
		// 削除
		_, err = mysql.TxWrite(tx, "DELETE FROM event_documents WHERE event_id = ?", id)
		if err != nil {
			tx.Rollback()
			return
		}
		updated.Documents = nil

		if *p.Documents.Slice != nil && len(**p.Documents.Slice) != 0 {
			// クエリを作成
			queryStr4 := "INSERT INTO event_documents (event_id, name, url) VALUES "
			var queryParams4 []interface{}
			for _, v := range **p.Documents.Slice {
				queryStr4 += "(?, ?),"
				queryParams4 = append(queryParams4, id, v.Name, v.Url)
			}
			queryStr4 = strings.TrimSuffix(queryStr4, ",")
			// 書込
			var result2 sql.Result
			result2, err = mysql.TxWrite(tx, queryStr4, queryParams4...)
			if err != nil {
				tx.Rollback()
				return
			}
			var eventDocumentId int64
			eventDocumentId, err = result2.LastInsertId()
			if err != nil {
				return
			}
			eventDocuments := []EventDocument{}
			for i, d := range **p.Documents.Slice {
				eventDocuments = append(
					eventDocuments,
					EventDocument{
						Id:   uint64(eventDocumentId + int64(i)),
						Name: d.Name,
						Url:  d.Url,
					},
				)
			}
			updated.Documents = eventDocuments
		}
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	e = updated
	return
}
