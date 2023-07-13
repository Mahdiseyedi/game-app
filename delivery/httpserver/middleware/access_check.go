package middleware

import (
	"fmt"
	"game-app/entity/permission"
	"game-app/pkg/claim"
	"game-app/pkg/errmsg"
	"game-app/service/authorizeservice"
	"github.com/labstack/echo/v4"
	"net/http"
)

func AccessCheck(service authorizeservice.Service, permissions ...permission.PermissionTitle) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			claims := claim.GetClaimsFromEchoContext(c)
			isAllowed, err := service.CheckAccess(claims.UserID, claims.Role, permissions...)

			if err != nil {
				fmt.Println("access control err", isAllowed, err)
				return c.JSON(http.StatusInternalServerError, echo.Map{
					"message": errmsg.ErrorMsgSomethingWentWrong,
				})
			}

			if !isAllowed {
				return c.JSON(http.StatusForbidden, echo.Map{
					"message": errmsg.ErrorMsgUserNotAllowed,
				})
			}

			return next(c)
		}
	}
}
