package userhandler

import (
	"game-app/param"
	"game-app/pkg/claim"
	"game-app/pkg/httpmsg"
	"github.com/labstack/echo/v4"
	"net/http"
)

func (h Handler) userProfile(c echo.Context) error {
	claims := claim.GetClaimsFromEchoContext(c)

	resp, err := h.userSvc.Profile(c.Request().Context(), param.ProfileRequest{UserID: claims.UserID})
	if err != nil {
		msg, code := httpmsg.Error(err)
		return echo.NewHTTPError(code, msg)
	}

	return c.JSON(http.StatusOK, resp)
}
