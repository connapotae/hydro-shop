package servers

import (
	"github.com/connapotae/hydro-shop/modules/appinfo/appinfoHandlers"
	"github.com/connapotae/hydro-shop/modules/appinfo/appinfoRepositories"
	"github.com/connapotae/hydro-shop/modules/appinfo/appinfoUsecases"
	"github.com/connapotae/hydro-shop/modules/middlewares/middlewaresHandlers"
	"github.com/connapotae/hydro-shop/modules/middlewares/middlewaresRepositories"
	"github.com/connapotae/hydro-shop/modules/middlewares/middlewaresUsecases"
	"github.com/connapotae/hydro-shop/modules/monitor/monitorHandlers"
	"github.com/connapotae/hydro-shop/modules/users/usersHandlers"
	"github.com/connapotae/hydro-shop/modules/users/usersRepositories"
	"github.com/connapotae/hydro-shop/modules/users/usersUsecases"
	"github.com/gofiber/fiber/v2"
)

type IModuleFactory interface {
	MonitorModule()
	UsersModule()
	AppinfoModule()
}

type moduleFactory struct {
	router fiber.Router
	server *server
	mid    middlewaresHandlers.IMiddlewaresHandler
}

func InitModule(r fiber.Router, s *server, mid middlewaresHandlers.IMiddlewaresHandler) IModuleFactory {
	return &moduleFactory{
		router: r,
		server: s,
		mid:    mid,
	}
}

func InitMiddlewares(s *server) middlewaresHandlers.IMiddlewaresHandler {
	repository := middlewaresRepositories.MiddlewaresRepository(s.db)
	usecase := middlewaresUsecases.MiddlewaresUsecase(repository)
	return middlewaresHandlers.MiddlewaresHandler(s.cfg, usecase)
}

func (m *moduleFactory) MonitorModule() {
	handler := monitorHandlers.MonitorHandler(m.server.cfg)

	m.router.Get("/", handler.HealthCheck)
}

func (m *moduleFactory) UsersModule() {
	repository := usersRepositories.UsersRepository(m.server.db)
	usecase := usersUsecases.UsersUsecase(m.server.cfg, repository)
	handler := usersHandlers.UsersHandler(m.server.cfg, usecase)

	router := m.router.Group("/users")

	router.Post("/signup", m.mid.ApiKeyAuth(), handler.SignUpCustomer)
	router.Post("/signin", m.mid.ApiKeyAuth(), handler.SignIn)
	router.Post("/refresh", m.mid.ApiKeyAuth(), handler.RefreshPassport)
	router.Post("/signout", m.mid.ApiKeyAuth(), handler.SignOut)
	router.Post("/signup-admin", m.mid.JwtAuth(), m.mid.Authorize(2), handler.SignUpAdmin)

	router.Get("/:user_id", m.mid.JwtAuth(), m.mid.ParamsCheck(), handler.GetUserProfile)
	router.Get("admin/token", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateAdminToken)

}

func (m *moduleFactory) AppinfoModule() {
	repository := appinfoRepositories.AppinfoRepository(m.server.db)
	usecase := appinfoUsecases.AppinfoUsecase(repository)
	handler := appinfoHandlers.AppinfoHandler(m.server.cfg, usecase)

	router := m.router.Group("/appinfo")

	// router.Get("/apikey", handler.GenerateApiKey)
	router.Get("/apikey", m.mid.JwtAuth(), m.mid.Authorize(2), handler.GenerateApiKey)
	router.Get("/categories", m.mid.ApiKeyAuth(), handler.FindCategory)
	router.Post("/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.AddCategory)
	router.Delete("/:category_id/categories", m.mid.JwtAuth(), m.mid.Authorize(2), handler.RemoveCateogory)
}
