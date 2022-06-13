package users

import (
	"errors"
	"prc_hub-api/mysql"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type PatchBody struct {
	Name                *string                   `json:"name" form:"name" validate:"omitempty,gte=1"`
	Email               *string                   `json:"email" form:"email" validate:"omitempty,email"`
	Password            *string                   `json:"password" form:"password" validate:"omitempty,gte=8"`
	GithubUsername      mysql.PatchNullJSONString `json:"github_username" validate:"omitempty,gte=1"`
	TwitterId           mysql.PatchNullJSONString `json:"twiiter_id" validate:"omitempty,gte=1"`
	PostEventAvailabled *bool                     `json:"post_event_availabled" validate:"omitempty"`
	Admin               *bool                     `json:"admin" validate:"omitempty"`
}

func (p *PatchBody) Validate() (err error) {
	if p.Name == nil &&
		p.Email == nil &&
		p.Password == nil &&
		p.GithubUsername.String == nil &&
		p.TwitterId.String == nil &&
		p.PostEventAvailabled == nil &&
		p.Admin == nil {
		err = errors.New("no update")
	}
	return
}

func Patch(id uint64, new PatchBody) (u User, invalidEmail bool, usedEmail bool, notFound bool, err error) {
	// userを取得
	u, notFound, err = GetById(id)
	if err != nil {
		return
	}
	if notFound {
		return
	}

	// クエリを作成
	queryStr := "UPDATE users SET "
	var queryParams []interface{}
	updated := u
	if new.Name != nil {
		queryStr += " name = ?,"
		queryParams = append(queryParams, new.Name)
		updated.Name = *new.Name
	}
	if new.Email != nil {
		queryStr += " email = ?,"
		queryParams = append(queryParams, new.Email)
		updated.Email = *new.Email
	}
	if new.Password != nil {
		queryStr += " password = ?"
		// Create password hash
		var hashed []byte
		hashed, err = bcrypt.GenerateFromPassword([]byte(*new.Password), 10)
		if err != nil {
			return
		}
		queryParams = append(queryParams, hashed)
		updated.Password = hashed
	}
	if new.GithubUsername.String != nil {
		if *new.GithubUsername.String != nil {
			queryStr += " github_username = ?,"
			queryParams = append(queryParams, **new.GithubUsername.String)
			updated.GithubUsername = *new.GithubUsername.String
		} else {
			queryStr += " github_username = ?,"
			queryParams = append(queryParams, nil)
			updated.GithubUsername = nil
		}
	}
	if new.TwitterId.String != nil {
		if *new.TwitterId.String != nil {
			queryStr += " twitter_id = ?,"
			queryParams = append(queryParams, **new.TwitterId.String)
			updated.TwitterId = *new.TwitterId.String
		} else {
			queryStr += " twitter_id = ?,"
			queryParams = append(queryParams, nil)
			updated.TwitterId = nil
		}
	}
	if new.PostEventAvailabled != nil {
		queryStr += " post_event_availabled = ?,"
		queryParams = append(queryParams, new.PostEventAvailabled)
		updated.PostEventAvailabled = *new.PostEventAvailabled
	}
	if new.Admin != nil {
		queryStr += " admin = ?,"
		queryParams = append(queryParams, new.Admin)
		updated.Admin = *new.Admin
	}
	queryStr = strings.TrimRight(queryStr, ",")
	queryStr += " WHERE id = ?"
	queryParams = append(queryParams, id)

	// 更新
	_, err = mysql.Write(queryStr, queryParams...)
	if err != nil {
		return
	}

	u = updated
	return
}
