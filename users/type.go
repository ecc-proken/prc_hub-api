package users

type User struct {
	Id                  uint64  `json:"id"`
	Name                string  `json:"name"`
	Email               string  `json:"email"`
	Password            []byte  `json:"-"`
	GithubUsername      *string `json:"github_username,omitempty"`
	TwitterId           *string `json:"twiiter_id,omitempty"`
	PostEventAvailabled bool    `json:"post_event_availabled"`
	Admin               bool    `json:"admin"`
}

type UserEmbed struct {
	Id             uint64  `json:"id"`
	Name           string  `json:"name"`
	GithubUsername *string `json:"github_username,omitempty"`
	TwitterId      *string `json:"twiiter_id,omitempty"`
}
