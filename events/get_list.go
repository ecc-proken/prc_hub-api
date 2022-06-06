package events

import (
	"prc_hub-api/mysql"
	"strings"
)

type GetQuery struct {
	Published       *bool   `json:"published" validate:"omitempty"`
	Title           *string `json:"name" validate:"omitempty"`
	TitleContain    *string `json:"name_contain" validate:"omitempty"`
	Location        *string `json:"location" validate:"omitempty"`
	LocationContain *string `json:"location_contain" validate:"omitempty"`
}

func Get(query GetQuery, userId *uint64, admin bool) (events []Event, err error) {
	// クエリを作成
	queryStr :=
		`SELECT
			e.id, e.user_id, e.title, e.location, e.published, e.completed, e.auto_notify_documents_enabled,
			doc.id, doc.name, doc.url
		FROM events e
		LEFT JOIN event_documents doc ON e.id = doc.event_id
		WHERE`
	queryParams := []interface{}{}

	if userId != nil && !admin {
		queryStr += " e.published = true OR e.user_id = ? AND"
		queryParams = append(queryParams, userId)
	}
	if userId == nil {
		queryStr += " e.published = true AND"
	}
	if userId != nil && query.Published != nil {
		queryStr += " e.published = ? AND"
		queryParams = append(queryParams, query.Published)
	}
	if query.Title != nil {
		queryStr += " e.title = ? AND"
		queryParams = append(queryParams, query.Title)
	}
	if query.TitleContain != nil {
		queryStr += " e.title LIKE ? AND"
		queryParams = append(queryParams, "%"+*query.TitleContain+"%")
	}
	if query.Location != nil {
		queryStr += " e.location = ? AND"
		queryParams = append(queryParams, query.Location)
	}
	if query.LocationContain != nil {
		queryStr += " e.location LIKE ?"
		queryParams = append(queryParams, "%"+*query.LocationContain+"%")
	}

	queryStr = strings.TrimSuffix(queryStr, "WHERE")
	queryStr = strings.TrimSuffix(queryStr, "AND")

	rows, err := mysql.Read(queryStr, queryParams...)
	if err != nil {
		return
	}
	defer rows.Close()

	// 読込中Event
	var tmpEvent *Event
	// 1行ずつ読込
	for rows.Next() {
		// 読込用変数
		e := Event{}
		var (
			tmpDocId   *uint64
			tmpDocName *string
			tmpDocUrl  *string
		)
		// 変数に割り当て
		err = rows.Scan(
			&e.Id, &e.UserId, &e.Title, &e.Location, &e.Published, &e.Completed, &e.AutoNotifyDocuments,
			&tmpDocId, &tmpDocName, &tmpDocUrl,
		)
		if err != nil {
			return
		}

		// 読込中のEventを更新
		if tmpEvent == nil {
			// 読込中のEventがない場合(初回に実行)
			// 新しく読み込んだEventを保持
			tmpEvent = &e
		} else if tmpEvent.Id != e.Id {
			// Eventが変わった場合
			// レスポンス用の配列に追加
			events = append(events, *tmpEvent)
			// 新しく読み込んだEventを保持
			tmpEvent = &e
		}
		// EventDocumentを追加
		if tmpDocId != nil && tmpDocName != nil && tmpDocUrl != nil {
			// 読込中のEventとIdが一致した場合
			// EventのDocumentを追加
			tmpEvent.Documents = append(tmpEvent.Documents, EventDocument{*tmpDocId, *tmpDocName, *tmpDocUrl})
		}
	}
	events = append(events, *tmpEvent)

	return
}
