package events

import (
	"net/http"
	"prc_hub-api/events"
	"prc_hub-api/flags"
	"prc_hub-api/jwt"
	"prc_hub-api/users"
	"strconv"

	jwtGo "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func Unparticipate(c echo.Context) (err error) {
	// jwtトークン確認
	t := c.Get("user").(*jwtGo.Token)
	claims, err := jwt.CheckToken(*flags.Get().JwtIssuer, t)
	if err != nil {
		c.Logger().Debug("401: " + err.Error())
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// userを取得
	u, notFound, err := users.GetById(claims.Id)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		c.Logger().Debug("404: user not found")
		return echo.ErrNotFound
	}

	// id
	idStr := c.Param("id")
	// string -> uint64
	_, err = strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// 404: Not found
		c.Logger().Debug("404: cannot parse `:id`")
		return echo.ErrNotFound
	}

	// datetime id
	datetimeIdStr := c.Param("dt_id")
	// string -> uint64
	datetimeId, err := strconv.ParseUint(datetimeIdStr, 10, 64)
	if err != nil {
		// 404: Not found
		c.Logger().Debug("404: cannot parse `:dt_id`")
		return echo.ErrNotFound
	}

	// 参加情報を削除
	notFound, err = events.Unparticipate(datetimeId, u.Id)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("404: event participate not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "event participate not found"}, "	")
	}

	// 204: No content
	c.Logger().Debug("204: unparticipate successful")
	return c.JSONPretty(http.StatusNoContent, map[string]string{"message": "Deleted"}, "	")
}
