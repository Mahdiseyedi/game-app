package backofficeuserhandler

import (
	"game-app/pkg/httpmsg"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) listUsers(c echo.Context) error {
	list, err := h.backofficeUserSvc.ListAllUsers()
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": list,
	})
}
