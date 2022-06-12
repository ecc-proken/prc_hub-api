package events

import "prc_hub-api/mysql"

func Participate(eventId uint64, datetimeId uint64, userId uint64) (e EventParticipate, eventNotFound bool, datetimeNotFound bool, userNotFound bool, err error) {
	// トランザクション開始
	tx, err := mysql.Begin()
	if err != nil {
		return
	}

	// userを取得
	r1, err := mysql.TxRead(tx, "SELECT id FROM users WHERE id = ?", userId)
	if err != nil {
		return
	}
	defer r1.Close()
	if !r1.Next() {
		userNotFound = true
		return
	}

	// eventを取得
	r2, err := mysql.TxRead(tx, "SELECT id FROM events WHERE id = ?", eventId)
	if err != nil {
		return
	}
	defer r2.Close()
	if !r2.Next() {
		eventNotFound = true
		return
	}

	// event_datetimeを取得
	r3, err := mysql.TxRead(tx, "SELECT id FROM event_datetimes WHERE id = ?", eventId)
	if err != nil {
		return
	}
	defer r3.Close()
	if !r3.Next() {
		datetimeNotFound = true
		return
	}

	// 参加登録情報を取得
	r4, err := mysql.TxRead(tx, "SELECT event_datetime_id FROM event_participates WHERE event_datetime_id = ? AND user_id = ?", datetimeId, userId)
	if err != nil {
		return
	}
	defer r4.Close()
	if r4.Next() {
		// 登録済みの場合は処理をスキップ
		return
	}

	// 書込
	_, err = mysql.TxWrite(tx, "INSERT INTO event_participates (event_datetime_id, user_id) VALUES (?, ?)", datetimeId, userId)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	e.EventDatetimeId = datetimeId
	e.UserId = userId
	return
}
