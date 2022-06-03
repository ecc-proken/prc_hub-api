package oauth_provider

import (
	"errors"
	"net/http"
	"prc_hub-api/flags"
	"prc_hub-api/jwt"
	"prc_hub-api/oauth2"
	"prc_hub-api/oauth2/github"

	jwtGo "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func Delete(c echo.Context) (err error) {
	// jwtトークン確認
	token := c.Get("user").(*jwtGo.Token)
	claims, err := jwt.CheckToken(*flags.Get().JwtIssuer, token)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}
	user_id := claims.Id

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
		c.Logger().Debug(errors.New("provider not found"))
		return echo.ErrNotFound
	}

	switch provider {
	case "github":
		// 書込
		notFound, err := github.Delete(user_id)
		if notFound {
			c.Logger().Debug("GitHub OAuth2 connection not found")
			return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "GitHub OAuth2 connection not found"}, "	")
		}
		if err != nil {
			c.Logger().Debug(err)
			return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
		}
	}

	// 204: No content
	return c.JSONPretty(http.StatusNoContent, map[string]string{"message": "Deleted"}, "	")
}
