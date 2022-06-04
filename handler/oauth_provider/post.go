package oauth_provider

import (
	"errors"
	"net/http"
	"prc_hub-api/flags"
	"prc_hub-api/jwt"
	"prc_hub-api/mysql"
	"prc_hub-api/oauth2"
	"prc_hub-api/oauth2/github"
	"prc_hub-api/users"

	jwtGo "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

type OAuth2GitHubPost struct {
	AccessToken string `json:"access_token" validate:"required"`
}

func Post(c echo.Context) (err error) {
	// jwtトークン確認
	token := c.Get("user").(*jwtGo.Token)
	claims, err := jwt.CheckToken(*flags.Get().JwtIssuer, token)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}
	user_id := claims.Id

	// Get user
	u, notFound, err := users.GetById(user_id)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "user not found"}, "	")
	}

	// Privider
	provider := c.Param("provider")
	switch provider {
	case oauth2.ProviderGitHub.String():
		if _, err = github.GetClient(); err != nil {
			// 404: Not found
			return echo.ErrNotFound
		}

	default:
		// 404: Not found
		c.Logger().Debug(errors.New("provider not found"))
		return echo.ErrNotFound
	}

	switch provider {
	case "github":
		// リクエストボディをバインド
		p := new(OAuth2GitHubPost)
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
		var a *github.Client
		a, err = github.GetClient()
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
		var o github.Owner
		o, err = a.GetOwner(p.AccessToken)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		// 書込
		_, err = github.Post(
			github.OAuth2{
				AccessToken: p.AccessToken,
				OwnerId:     o.Id,
			},
			user_id,
		)
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}

		// ユーザー情報の更新
		// github_usernameを追加
		tmpUsername := &o.Name
		_, _, _, notFound, err = users.Patch(user_id, users.PatchBody{GithubUsername: mysql.PatchNullJSONString{String: &tmpUsername}})
		if err != nil {
			return
		}
		if notFound {
			// ユーザー情報変更に失敗
			return c.JSONPretty(http.StatusConflict, map[string]string{"message": "connot update user, conflict found"}, "	")
		}
	}

	// トークンを生成
	t, err := jwt.GenerateToken(u, *flags.Get().JwtIssuer, *flags.Get().JwtSecret)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// Cookieを追加
	c.SetCookie(&http.Cookie{
		Name:     "token",
		Value:    t,
		HttpOnly: true,
	})

	// 200: Success
	return c.JSONPretty(http.StatusOK, map[string]string{"message": "Success"}, "	")
}
