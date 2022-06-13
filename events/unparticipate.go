package events

import (
	"prc_hub-api/mysql"
)

func Unparticipate(datetimeId uint64, userId uint64) (notFound bool, err error) {
	// トランザクション開始
	tx, err := mysql.Begin()
	if err != nil {
		return
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// 参加登録情報を取得
	r, err := tx.Query("SELECT true FROM event_participates WHERE event_datetime_id = ? AND user_id = ?", datetimeId, userId)
	if err != nil {
		return
	}
	defer r.Close()
	// 登録済みか確認
	if !r.Next() {
		notFound = true
		return
	}
	err = r.Close()
	if err != nil {
		return
	}

	// 削除
	_, err = tx.Exec("DELETE FROM event_participates WHERE event_datetime_id = ? AND user_id = ?", datetimeId, userId)
	if err != nil {
		return
	}

	err = tx.Commit()
	if err != nil {
		return
	}

	return
}
