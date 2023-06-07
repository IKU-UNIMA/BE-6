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

	fakultas := v1.Group("/fakultas", customMiddleware.Authentication)
	fakultas.GET("", handler.GetAllFakultasHandler)
	fakultas.GET("/:id", handler.GetFakultasByIdHandler)
	fakultas.POST("", handler.InsertFakultasHandler)
	fakultas.PUT("/:id", handler.EditFakultasHandler)
	fakultas.DELETE("/:id", handler.DeleteFakultasHandler)

	prodi := v1.Group("/prodi", customMiddleware.Authentication)
	prodi.GET("", handler.GetAllProdiHandler)
	prodi.GET("/:id", handler.GetProdiByIdHandler)
	prodi.POST("", handler.InsertProdiHandler)
	prodi.PUT("/:id", handler.EditProdiHandler)
	prodi.DELETE("/:id", handler.DeleteProdiHandler)

	akun := v1.Group("/akun")
	akun.POST("/login", handler.LoginHandler)
	akun.PATCH("/password/change", handler.ChangePasswordHandler, customMiddleware.Authentication)
	akun.PATCH("/password/reset/:id", handler.ResetPasswordHandler, customMiddleware.Authentication)

	admin := v1.Group("/admin", customMiddleware.Authentication)
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

	kerjasama := v1.Group("/kerjasama", customMiddleware.Authentication)
	kerjasama.PATCH("/:id/dokumen", handler.EditDokumenKerjasamaHandler)

	IA := kerjasama.Group("/IA")
	IA.GET("", handler.GetAllKerjasamaIAHandler)
	IA.GET("/:id", handler.GetKerjasamaIAByIdHandler)
	IA.GET("/dasar-kerjasama", handler.GetDasarKerjasamaIAHandler)
	IA.POST("", handler.InsertKerjasamaIAHandler, customMiddleware.GrantAdminIKU6)
	IA.PUT("/:id", handler.EditKerjasamaIAHandler, customMiddleware.GrantAdminIKU6)
	IA.DELETE("/:id", handler.DeleteKerjasamaIAHandler, customMiddleware.GrantAdminIKU6)

	MOA := kerjasama.Group("/MOA")
	MOA.GET("", handler.GetAllKerjasamaMOAHandler)
	MOA.GET("/:id", handler.GetKerjasamaMOAByIdHandler)
	MOA.GET("/dasar-kerjasama", handler.GetDasarKerjasamaMOAHandler)
	MOA.POST("", handler.InsertKerjasamaMOAHandler, customMiddleware.GrantAdminIKU6)
	MOA.PUT("/:id", handler.EditKerjasamaMOAHandler, customMiddleware.GrantAdminIKU6)
	MOA.DELETE("/:id", handler.DeleteKerjasamaMOAHandler, customMiddleware.GrantAdminIKU6)

	MOU := kerjasama.Group("/MOU")
	MOU.GET("", handler.GetAllKerjasamaMOUHandler)
	MOU.GET("/:id", handler.GetKerjasamaMOUByIdHandler)
	MOU.POST("", handler.InsertKerjasamaMOUHandler, customMiddleware.GrantAdminIKU6)
	MOU.PUT("/:id", handler.EditKerjasamaMOUHandler, customMiddleware.GrantAdminIKU6)
	MOU.DELETE("/:id", handler.DeleteKerjasamaMOUHandler, customMiddleware.GrantAdminIKU6)

	dashboard := v1.Group("/dashboard", customMiddleware.Authentication)
	dashboard.GET("/tahun/:tahun", handler.GetDashboardHandler, customMiddleware.GrantAdminIKU6AndRektor)
	dashboard.GET("/fakultas/:fakultas/:tahun", handler.GetDashboardByFakultasHandler, customMiddleware.GrantAdminIKU6AndRektor)
	dashboard.PATCH("/target", handler.InsertTargetHandler, customMiddleware.GrantAdminIKU6)

	return app

}
