package users

import (
	"net/http"
	"prc_hub-api/flags"
	"prc_hub-api/jwt"
	"prc_hub-api/users"
	"strconv"

	jwtGo "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func DeleteById(c echo.Context) (err error) {
	// jwtトークン確認
	u := c.Get("user").(*jwtGo.Token)
	claims, err := jwt.CheckToken(*flags.Get().JwtIssuer, u)
	if err != nil {
		c.Logger().Debug("401: " + err.Error())
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// id
	idStr := c.Param("id")
	// string -> uint64
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// 404: Not found
		c.Logger().Debug("404: cannot parse `:id`")
		return echo.ErrNotFound
	}

	// 権限確認
	if !claims.Admin && claims.Id != id {
		// Admin権限なし 且つ IDが自分ではない
		// 403: Forbidden
		c.Logger().Debug("403: you cannot delete this user")
		return c.JSONPretty(http.StatusForbidden, map[string]string{"message": "you cannot delete this user"}, "	")
	}

	if claims.Id == id && claims.Email == *flags.Get().AdminEmail {
		// name=adminは削除不可
		// 405: Method not allowed
		c.Logger().Debug("405: cannot delete admin user")
		return c.JSONPretty(http.StatusMethodNotAllowed, map[string]string{"message": "cannot delete admin user"}, "	")
	}

	// 削除
	notFound, err := users.Delete(id)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("404: user not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "user not found"}, "	")
	}

	// 204: No content
	c.Logger().Debug("204: delete user successful")
	return c.JSONPretty(http.StatusNoContent, map[string]string{"message": "Deleted"}, "	")
}
