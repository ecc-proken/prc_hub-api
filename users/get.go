package users

import "prc_hub-api/mysql"

func GetById(id uint64) (u User, notFound bool, err error) {
	// 読込
	rows, err := mysql.Read(
		`SELECT name, email, password, github_username, twitter_id, post_event_availabled, admin
		FROM users WHERE id = ?`,
		id,
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

	err = rows.Scan(&u.Name, &u.Email, &u.Password, &u.GithubUsername, &u.TwitterId, &u.PostEventAvailabled, &u.Admin)
	if err != nil {
		return
	}

	u.Id = id
	return
}

func GetByEmail(email string) (u User, notFound bool, err error) {
	// 読込
	rows, err := mysql.Read(
		`SELECT id, name, password, github_username, twitter_id, post_event_availabled, admin
		FROM users WHERE email = ?`,
		email,
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
	err = rows.Scan(&u.Id, &u.Name, &u.Password, &u.GithubUsername, &u.TwitterId, &u.PostEventAvailabled, &u.Admin)
	if err != nil {
		return
	}

	u.Email = email
	return
}
