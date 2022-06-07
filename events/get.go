package events

import (
	"prc_hub-api/mysql"
)

func GetById(id uint64) (events []Event, err error) {
	rows, err := mysql.Read(
		`SELECT
			e.id, e.user_id, e.title, e.location, e.published, e.completed, e.auto_notify_documents_enabled,
			doc.id, doc.name, doc.url
		FROM events e
		LEFT JOIN event_documents doc WHERE e.id = doc.event_id
		WHERE e.id = ?`,
		id,
	)
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
		} else if tmpDocId != nil && tmpDocName != nil && tmpDocUrl != nil {
			// 読込中のEventとIdが一致した場合
			// EventのDocumentを追加
			tmpEvent.Documents = append(tmpEvent.Documents, EventDocument{*tmpDocId, *tmpDocName, *tmpDocUrl})
		}
	}

	return
}
