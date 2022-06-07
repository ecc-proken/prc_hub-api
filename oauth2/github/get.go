package github

import (
	"prc_hub-api/mysql"
)

func Get(user_id uint64) (o OAuth2, notFound bool, err error) {
	// 読込
	rows, err := mysql.Read(
		"SELECT access_token, owner_id FROM github_oauth2_tokens WHERE user_id = ?",
		user_id,
	)
	if err != nil {
		return
	}
	defer rows.Close()

	if !rows.Next() {
		// Not found
		notFound = true
		return
	}
	err = rows.Scan(&o.AccessToken, &o.OwnerId)
	if err != nil {
		return
	}

	return
}
