package route

import (
	"BE-6/src/api/handler"
	"BE-6/src/util/validation"

	customMiddleware "BE-6/src/api/middleware"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func InitServer() *echo.Echo {

	app := echo.New()
	app.Use(middleware.CORS())

	app.Validator = &validation.CustomValidator{Validator: validator.New()}

	app.GET("", func(c echo.Context) error {
		return c.JSON(200, "Welcome to IKU 6 API")
	})

	v1 := app.Group("/api/v1")

	akun := v1.Group("/akun")
	akun.POST("/login", handler.LoginHandler)
	akun.PATCH("/password/change", handler.ChangePasswordHandler, customMiddleware.Authentication)
	akun.PATCH("/password/reset/:id", handler.ResetPasswordHandler, customMiddleware.Authentication, customMiddleware.GrantAdminUmum)

	admin := v1.Group("/admin", customMiddleware.Authentication, customMiddleware.GrantAdminUmum)
	admin.GET("", handler.GetAllAdminHandler)
	admin.GET("/:id", handler.GetAdminByIdHandler)
	admin.POST("", handler.InsertAdminHandler)
	admin.PUT("/:id", handler.EditAdminHandler)
	admin.DELETE("/:id", handler.DeleteAdminHandler)

	rektor := v1.Group("/rektor", customMiddleware.Authentication)
	rektor.GET("", handler.GetAllRektorHandler)
	rektor.GET("/:id", handler.GetRektorByIdHandler)
	rektor.POST("", handler.InsertRektorHandler)
	rektor.PUT("/:id", handler.EditRektorHandler)
	rektor.DELETE("/:id", handler.DeleteRektorHandler)

	kerjsamaIA := v1.Group("/kerjasamaIA", customMiddleware.Authentication)
	kerjsamaIA.GET("", handler.GetAllKerjasamaIAHandler)
	kerjsamaIA.GET("/:id", handler.GetKerjasamaIAByIdHandler)
	kerjsamaIA.POST("", handler.InsertKerjasamaIAHandler, customMiddleware.GrantAdminIKU6)
	kerjsamaIA.PUT("/:id", handler.EditKerjasamaIAHandler, customMiddleware.GrantAdminIKU6)
	kerjsamaIA.DELETE("/:id", handler.DeleteKerjasamaIAHandler, customMiddleware.GrantAdminIKU6)

	kerjsamaMOA := v1.Group("/kerjasamaMOA", customMiddleware.Authentication, customMiddleware.GrantAdminUmum)
	kerjsamaMOA.GET("", handler.GetAllKerjasamaMOAHandler)
	kerjsamaMOA.GET("/:id", handler.GetKerjasamaMOAByIdHandler)
	kerjsamaMOA.POST("", handler.InsertKerjasamaMOAHandler, customMiddleware.GrantAdminIKU6)
	kerjsamaMOA.PUT("/:id", handler.EditKerjasamaMOAHandler, customMiddleware.GrantAdminIKU6)
	kerjsamaMOA.DELETE("/:id", handler.DeleteKerjasamaMOAHandler, customMiddleware.GrantAdminIKU6)

	kerjsamaMOU := v1.Group("/kerjasamaMOU", customMiddleware.Authentication, customMiddleware.GrantAdminUmum)
	kerjsamaIA.GET("", handler.GetAllKerjasamaMOUHandler)
	kerjsamaMOU.GET("/:id", handler.GetKerjasamaMOAByIdHandler)
	kerjsamaMOU.POST("", handler.InsertKerjasamaMOUHandler, customMiddleware.GrantAdminIKU6)
	kerjsamaMOU.PUT("/:id", handler.EditKerjasamaMOUHandler, customMiddleware.GrantAdminIKU6)
	kerjsamaMOU.DELETE("/:id", handler.DeleteKerjasamaMOUHandler, customMiddleware.GrantAdminIKU6)
	return app

}
