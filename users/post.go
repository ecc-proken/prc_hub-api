package users

import (
	"prc_hub-api/mysql"

	"golang.org/x/crypto/bcrypt"
)

type PostBody struct {
	Name           string  `json:"name" form:"name" validate:"required"`
	Email          string  `json:"email" form:"email" validate:"required,email"`
	Password       string  `json:"password" form:"password" validate:"required"`
	GithubUsername *string `json:"github_username" validate:"omitempty"`
	TwitterId      *string `json:"twiiter_id" validate:"omitempty"`
}

func Post(post PostBody) (u User, invalidEmail bool, usedEmail bool, err error) {
	// メールアドレスがすでに使用されていないか確認
	_, notFound, err := GetByEmail(post.Email)
	if err != nil {
		return
	}
	if !notFound {
		// 使用済みメールアドレス
		invalidEmail = true
		return
	}

	// パスワードをハッシュ化
	hashed, err := bcrypt.GenerateFromPassword([]byte(post.Password), 10)
	if err != nil {
		return
	}

	// 書込
	result, err := mysql.Write(
		`INSERT INTO users (name, email, password, github_username, twitter_id)
			VALUES (?, ?, ?, ?, ?)`,
		post.Name,
		post.Email,
		hashed,
		post.GithubUsername,
		post.TwitterId,
	)
	if err != nil {
		return
	}
	// Insertした行のIdを取得
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	u.Id = uint64(id)
	u.Name = post.Name
	u.Email = post.Email
	u.GithubUsername = post.GithubUsername
	u.TwitterId = post.TwitterId
	return
}

func PostAdmin(post PostBody) (u User, invalidEmail bool, usedEmail bool, err error) {
	// メールアドレスがすでに使用されていないか確認
	_, notFound, err := GetByEmail(post.Email)
	if err != nil {
		return
	}
	if !notFound {
		// 使用済みメールアドレス
		invalidEmail = true
		return
	}

	// パスワードをハッシュ化
	hashed, err := bcrypt.GenerateFromPassword([]byte(post.Password), 10)
	if err != nil {
		return
	}

	// 書込
	result, err := mysql.Write(
		`INSERT INTO users (name, email, password, github_username, twitter_id, post_event_availabled, admin)
			VALUES (?, ?, ?, ?, ?, ?, ?)`,
		post.Name,
		post.Email,
		hashed,
		post.GithubUsername,
		post.TwitterId,
		true,
		true,
	)
	if err != nil {
		return
	}
	// Insertした行のIdを取得
	id, err := result.LastInsertId()
	if err != nil {
		return
	}

	u.Id = uint64(id)
	u.Name = post.Name
	u.Email = post.Email
	u.GithubUsername = post.GithubUsername
	u.TwitterId = post.TwitterId
	u.PostEventAvailabled = true
	u.Admin = true
	return
}
