package events

import "prc_hub-api/mysql"

func Unparticipate(datetimeId uint64, userId uint64) (notFound bool, err error) {
	// 参加登録情報を取得
	r, err := mysql.Read("SELECT true FROM event_participates WHERE event_datetime_id = ? AND user_id = ?", datetimeId, userId)
	if err != nil {
		return
	}
	defer r.Close()

	// 登録済みか確認
	found := false
	if r.Next() {
		err = r.Scan(&found)
		if err != nil {
			return
		}
	}
	if !found {
		notFound = true
		return
	}

	// 削除
	_, err = mysql.Write("DELETE FROM event_participates WHERE datetime_id = ? AND user_id = ?", datetimeId, userId)
	return
}
