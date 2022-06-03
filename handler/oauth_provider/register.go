package oauth_provider

import (
	"encoding/json"
	"net/http"
	"prc_hub-api/flags"
	"prc_hub-api/jwt"
	"prc_hub-api/oauth2"
	"prc_hub-api/oauth2/github"
	"prc_hub-api/users"

	"github.com/labstack/echo/v4"
)

type UserPostOverGitHubOAuth2 struct {
	AccessToken string `json:"access_token" validate:"required"`
	Password    string `json:"password" validate:"required"`
}

func Register(c echo.Context) (err error) {
	// Privider
	provider := c.Param("provider")
	switch provider {
	case oauth2.ProviderGitHub.String():
		if *flags.Get().GithubClientId == "" || *flags.Get().GithubClientSecret == "" {
			// 404: Not found
			return echo.ErrNotFound
		}

	default:
		// 404: Not found
		c.Logger().Debug(err)
		return echo.ErrNotFound
	}

	// GitHubの登録情報を取得
	var u users.User
	var name string
	var email string

	switch provider {
	case "github":
		// リクエストボディをバインド
		p := new(UserPostOverGitHubOAuth2)
		if err = c.Bind(p); err != nil {
			// 400: Bad request
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
		}

		// リクエストボディを検証
		if err = c.Validate(p); err != nil {
			// 422: Unprocessable entity
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
		}

		// GitHubの登録情報を取得
		a, err := github.New(*flags.Get().GithubClientId, *flags.Get().GithubClientSecret)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		// 登録情報の取得
		o, err := a.GetOwner(p.AccessToken)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		name = o.Name
		// primaryEmailの取得
		e, err := a.GetOwnerPrimaryEmail(p.AccessToken)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		email = e.Email

		// 書込
		u, invalidEmail, usedEmail, err := users.Post(users.PostBody{Name: name, Email: email, Password: p.Password, GithubUsername: &name})
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		if invalidEmail {
			// 422: Unprocessable entity
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
		}
		if usedEmail {
			// 409: Conflict
			c.Logger().Debug("email already used")
			return c.JSONPretty(http.StatusConflict, map[string]string{"message": "email already used"}, "	")
		}

		// 書込
		_, err = github.Post(
			github.OAuth2{
				AccessToken: p.AccessToken,
				OwnerId:     o.Id,
			},
			u.Id,
		)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
	}

	// トークンを生成
	t, err := jwt.GenerateToken(u, *flags.Get().JwtIssuer, *flags.Get().JwtSecret)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// jsonにjwtトークンを追加
	b, err := json.Marshal(u)
	if err != nil {
		return
	}
	m := map[string]interface{}{"token": t}
	err = json.Unmarshal(b, &m)
	if err != nil {
		return
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, m, "	")
}
