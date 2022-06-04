package github

import "prc_hub-api/mysql"

func Post(post OAuth2, user_id uint64) (o OAuth2, err error) {
	// 読込
	_, notFound, err := Get(user_id)
	if err != nil {
		return
	}
	if !notFound {
		// 古いデータを削除
		_, err = Delete(user_id)
		if err != nil {
			return
		}
	}

	// 書込
	_, err = mysql.Write(
		"INSERT INTO github_oauth2_tokens (user_id, access_token, owner_id) VALUES(?, ?, ?)",
		user_id, post.AccessToken, post.OwnerId,
	)
	if err != nil {
		return
	}

	o.AccessToken = post.AccessToken
	o.OwnerId = post.OwnerId
	return
}
