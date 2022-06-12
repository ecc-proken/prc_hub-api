package events

import "prc_hub-api/mysql"

func Participate(datetimeId uint64, userId uint64) (e EventParticipate, notFound bool, err error) {
	// 参加登録情報を取得
	r, err := mysql.Read("SELECT true FROM event_participates WHERE event_datetime_id = ? AND user_id = ?", datetimeId, userId)
	if err != nil {
		return
	}
	defer r.Close()

	// 登録済みの場合は処理をスキップ
	exists := false
	if r.Next() {
		err = r.Scan(&exists)
		if err != nil {
			return
		}
	}
	if exists {
		return
	}

	// 書込
	_, err = mysql.Write("INSERT INTO event_participates (event_datetime_id, user_id) VALUES (?, ?)", datetimeId, userId)

	e.EventDatetimeId = datetimeId
	e.UserId = userId
	return
}
