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

func Participate(c echo.Context) (err error) {
	// jwtトークン確認
	t := c.Get("user").(*jwtGo.Token)
	claims, err := jwt.CheckToken(*flags.Get().JwtIssuer, t)
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

	// datetime id
	datetimeIdStr := c.Param("dt_id")
	// string -> uint64
	datetimeId, err := strconv.ParseUint(datetimeIdStr, 10, 64)
	if err != nil {
		// 404: Not found
		c.Logger().Debug("404: cannot parse `:dt_id`")
		return echo.ErrNotFound
	}

	// 参加情報を登録
	ep, eventNotFound, datetimeNotFound, userNotFound, err := events.Participate(id, datetimeId, claims.Id)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if eventNotFound {
		// 404: Not found
		c.Logger().Debug("404: event not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "event not found"}, "	")
	}
	if datetimeNotFound {
		// 404: Not found
		c.Logger().Debug("404: event datetime not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "event datetime not found"}, "	")
	}
	if userNotFound {
		// 404: Not found
		c.Logger().Debug("404: user not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "user not found"}, "	")
	}

	// 200: Success
	c.Logger().Debug("200: participate  successful")
	return c.JSONPretty(http.StatusOK, ep, "	")
}
