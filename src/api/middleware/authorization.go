package middleware

import (
	"BE-6/src/util"
	"net/http"

	"github.com/labstack/echo/v4"
)

func GrantAdminUmum(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		if claims["role"].(string) != util.ADMIN ||
			claims["bagian"].(string) != util.UMUM {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}

func GrantAdminIKU6(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		if claims["role"].(string) != util.ADMIN ||
			claims["bagian"].(string) != util.IKU6 {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}

func GrantAdminIKU6AndRektor(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		role := claims["role"].(string)
		bagian := claims["bagian"].(string)
		if role != string(util.REKTOR) && role != string(util.ADMIN) {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		if role == string(util.ADMIN) && bagian != util.IKU6 {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}

func GrantDosen(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		claims := util.GetClaimsFromContext(c)
		if claims["role"].(string) != string(util.ADMIN) {
			return util.FailedResponse(http.StatusUnauthorized, nil)
		}

		return next(c)
	}
}
