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

func Participate(c echo.Context) (err error) {
	// jwtトークン確認
	t := c.Get("user").(*jwtGo.Token)
	claims, err := jwt.CheckToken(*flags.Get().JwtIssuer, t)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// userを取得
	u, notFound, err := users.GetById(claims.Id)
	if err != nil {
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		return echo.ErrNotFound
	}

	// id
	idStr := c.Param("id")
	// string -> uint64
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		// 404: Not found
		return echo.ErrNotFound
	}
	e, notFound, err := events.GetById(id)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		return echo.ErrNotFound
	}

	// datetime id
	datetimeIdStr := c.Param(":dt_id")
	// string -> uint64
	datetimeId, err := strconv.ParseUint(datetimeIdStr, 10, 64)
	if err != nil {
		// 404: Not found
		return echo.ErrNotFound
	}
	notFound = true
	for _, dt := range e.Datetimes {
		if datetimeId == dt.Id {
			notFound = false
		}
	}
	if notFound {
		// 404: Not found
		return echo.ErrNotFound
	}

	// 参加情報を登録
	ep, notFound, err := events.Participate(datetimeId, u.Id)
	if err != nil {
		c.Logger().Debug(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		return echo.ErrNotFound
	}

	// 200: Success
	return c.JSONPretty(http.StatusOK, ep, "	")
}
