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

func LineNotify(c echo.Context) (err error) {
	if flags.Get().LineNotifyToken == nil || *flags.Get().LineNotifyToken == "" {
		// 404: Not found
		c.Logger().Debug("404: LINE notify skipped")
		return echo.ErrNotFound
	}

	// jwtトークン確認
	u := c.Get("user").(*jwtGo.Token)
	_, err = jwt.CheckToken(*flags.Get().JwtIssuer, u)
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

	// eventを取得
	e, notFound, err := events.GetById(id)
	if err != nil {
		c.Logger().Fatal(err.Error())
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("404: event not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "event not found"}, "	")
	}

	// LINE notify
	err = e.NotifyLINE(*flags.Get().LineNotifyToken, *flags.Get().FrontUrl)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// 200: Success
	c.Logger().Debug("200: notify event to LINE successful")
	return c.JSONPretty(http.StatusOK, map[string]string{"message": "Success"}, "	")
}

func LineNotifyDocuments(c echo.Context) (err error) {
	if flags.Get().LineNotifyToken == nil || *flags.Get().LineNotifyToken == "" {
		// 404: Not found
		c.Logger().Debug("404: LINE notify skipped")
		return echo.ErrNotFound
	}

	// jwtトークン確認
	u := c.Get("user").(*jwtGo.Token)
	_, err = jwt.CheckToken(*flags.Get().JwtIssuer, u)
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

	// eventを取得
	e, notFound, err := events.GetById(id)
	if err != nil {
		c.Logger().Fatal(err.Error())
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("404: event not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "event not found"}, "	")
	}

	// LINE notify
	err = e.NotifyLINEDocuments(*flags.Get().LineNotifyToken)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}

	// 200: Success
	c.Logger().Debug("200: notify event to LINE successful")
	return c.JSONPretty(http.StatusOK, map[string]string{"message": "Success"}, "	")
}
