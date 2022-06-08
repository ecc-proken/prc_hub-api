package users

import (
	"prc_hub-api/mysql"
	"strings"
)

type GetQuery struct {
	Name                *string  `json:"name" validate:"omitempty"`
	NameContain         *string  `json:"name_contain" validate:"omitempty"`
	Email               *string  `json:"email" validate:"omitempty"`
	PostEventAvailabled *bool    `json:"post_event_availabled" validate:"omitempty"`
	Admin               *bool    `json:"admin" validate:"omitempty"`
	Ids                 []uint64 `json:"-" validate:"omitempty"`
}

func Get(query GetQuery) (users []User, err error) {
	// クエリを作成
	queryStr := "SELECT id, name, email, github_username, twitter_id, post_event_availabled, admin FROM users WHERE"
	queryParams := []interface{}{}

	if query.PostEventAvailabled != nil {
		queryStr += " post_event_availabled = ? AND"
		queryParams = append(queryParams, query.PostEventAvailabled)
	}
	if query.Admin != nil {
		queryStr += " admin = ? AND"
		queryParams = append(queryParams, query.Admin)
	}
	if query.Name != nil {
		queryStr += " name = ? AND"
		queryParams = append(queryParams, query.Name)
	}
	if query.NameContain != nil {
		queryStr += " name LIKE ? AND"
		queryParams = append(queryParams, "%"+*query.NameContain+"%")
	}
	if query.Email != nil {
		queryStr += " email = ? AND"
		queryParams = append(queryParams, query.Email)
	}
	if len(query.Ids) != 0 {
		queryStr += " id IN ("
		for _, id := range query.Ids {
			queryStr += " ?,"
			queryParams = append(queryParams, id)
		}
		queryStr = strings.TrimSuffix(queryStr, ",")
		queryStr += " ) AND"
	}
	queryStr = strings.TrimSuffix(queryStr, " WHERE")
	queryStr = strings.TrimSuffix(queryStr, " AND")

	// 読込
	rows, err := mysql.Read(queryStr, queryParams...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		u := User{}
		err = rows.Scan(&u.Id, &u.Name, &u.Email, &u.GithubUsername, &u.TwitterId, &u.PostEventAvailabled, &u.Admin)
		if err != nil {
			return
		}
		users = append(users, u)
	}

	return
}
