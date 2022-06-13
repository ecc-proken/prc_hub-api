package events

import (
	"net/http"
	"prc_hub-api/events"
	"prc_hub-api/flags"
	"prc_hub-api/jwt"
	"strconv"
	"strings"

	jwtGo "github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
)

func PatchById(c echo.Context) (err error) {
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
		c.Logger().Debug("403: you cannot update this event")
		return c.JSONPretty(http.StatusForbidden, map[string]string{"message": "you cannot update this event"}, "	")
	}

	// リクエストボディをバインド
	p := new(events.PatchBody)
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
	if err = p.Validate(); err != nil {
		// 400: Bad request
		c.Logger().Debug("400: " + err.Error())
		return c.JSONPretty(http.StatusBadRequest, map[string]string{"message": err.Error()}, "	")
	}

	// 更新
	e, notFound, notFoundUserIds, err := events.Patch(id, *p)
	if err != nil {
		c.Logger().Fatal(err)
		return c.JSONPretty(http.StatusInternalServerError, map[string]string{"message": err.Error()}, "	")
	}
	if notFound {
		// 404: Not found
		c.Logger().Debug("404: event not found")
		return c.JSONPretty(http.StatusNotFound, map[string]string{"message": "event not found"}, "	")
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
