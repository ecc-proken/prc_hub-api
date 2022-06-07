package users

import (
	"errors"
	"net/http"
	"prc_hub-api/flags"
	"prc_hub-api/jwt"
	"prc_hub-api/users"

	jwtGo "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func DeleteOwn(c echo.Context) (err error) {
	// jwtトークン確認
	u := c.Get("user").(*jwtGo.Token)
	claims, err := jwt.CheckToken(*flags.Get().JwtIssuer, u)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}
	id := claims.Id

	if claims.Admin && claims.Email == *flags.Get().AdminEmail {
		// name=adminは削除不可
		return c.JSONPretty(http.StatusMethodNotAllowed, map[string]string{"message": "cannot delete admin user"}, "	")
	}

	// 削除
	notFound, err := users.Delete(id)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug(errors.New("user not found"))
		return echo.ErrNotFound
	}

	// 204: No content
	return c.JSONPretty(http.StatusNoContent, map[string]string{"message": "Deleted"}, "	")
}
