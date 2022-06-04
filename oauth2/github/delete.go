package github

import "prc_hub-api/mysql"

func Delete(user_id uint64) (notFound bool, err error) {
	// 削除
	result, err := mysql.Write(
		"DELETE FROM github_oauth2_tokens WHERE user_id = ?",
		user_id,
	)
	if err != nil {
		return false, err
	}
	// Deleteの変更行数を取得
	affectedRowCount, err := result.RowsAffected()
	if err != nil {
		return false, err
	}
	notFound = affectedRowCount == 0

	return
}
