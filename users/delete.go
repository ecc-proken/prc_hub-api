package users

import "prc_hub-api/mysql"

func Delete(id uint64) (notFound bool, err error) {
	// 削除
	result, err := mysql.Write(
		"DELETE FROM users WHERE id = ?",
		id,
	)
	if err != nil {
		return false, err
	}
	// Deleteの変更行数を取得
	affectedRowCount, err := result.RowsAffected()
	if err != nil {
		return
	}
	notFound = affectedRowCount == 0

	return
}
