package events

import (
	"net/http"
	"prc_hub-api/events"
	"prc_hub-api/flags"
	"prc_hub-api/jwt"
	"prc_hub-api/users"
	"strconv"
	"strings"

	jwtGo "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func Post(c echo.Context) (err error) {
	// jwtトークン確認
	t := c.Get("user").(*jwtGo.Token)
	claims, err := jwt.CheckToken(*flags.Get().JwtIssuer, t)
	if err != nil {
		c.Logger().Debug("401: " + err.Error())
		return c.JSONPretty(http.StatusUnauthorized, map[string]string{"message": err.Error()}, "	")
	}

	// 権限確認
	u, notFound, err := users.GetById(claims.Id)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 403: Not found
		c.Logger().Debug("403: user not found")
		return echo.ErrForbidden
	}
	if !u.PostEventAvailabled {
		// 403: Forbidden
		c.Logger().Debug("403: you cannot create event")
		return c.JSONPretty(http.StatusForbidden, map[string]string{"message": "you cannot create event"}, "	")
	}

	// リクエストボディをバインド
	p := new(events.PostBody)
	if err = c.Bind(p); err != nil {
		// 400: Bad request
		c.Logger().Debug("400: " + err.Error())
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// リクエストボディを検証
	if err = c.Validate(p); err != nil {
		// 422: Unprocessable entity
		c.Logger().Debug("422: " + err.Error())
		return c.JSONPretty(http.StatusUnprocessableEntity, map[string]string{"message": err.Error()}, "	")
	}

	// 書込
	e, notFoundUserIds, err := events.Post(claims.Id, *p)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if len(notFoundUserIds) != 0 {
		msg := "user not found (id:"
		for _, id := range notFoundUserIds {
			msg += " " + strconv.FormatUint(id, 10) + ","
		}
		msg = strings.TrimSuffix(msg, ",")
		msg += ")"

		// 404: Not found
		c.Logger().Debug("404: " + msg)
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": msg}, "	")
	}

	// 200: Success
	c.Logger().Debug("200: patch event successful")
	return c.JSONPretty(http.StatusOK, e, "	")
}
