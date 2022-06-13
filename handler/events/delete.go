package events

import (
	"net/http"
	"prc_hub-api/events"
	"prc_hub-api/flags"
	"prc_hub-api/jwt"
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

	// Eventを取得
	e, notFound, err := events.GetById(id)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("404: event not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "event not found"}, "	")
	}
	if !claims.Admin && claims.Id != e.UserId {
		// 403: Forbidden
		c.Logger().Debug("403: you cannot delete this event")
		return c.JSONPretty(http.StatusForbidden, map[string]string{"message": "you cannot delete this event"}, "	")
	}

	// 削除
	notFound, err = events.Delete(id)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("404: event not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "event not found"}, "	")
	}

	// 204: No content
	c.Logger().Debug("204: event deleted")
	return c.JSONPretty(http.StatusNoContent, map[string]string{"message": "Deleted"}, "	")
}
