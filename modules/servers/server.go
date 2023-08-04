package servers

import (
	"encoding/json"
	"log"
	"os"
	"os/signal"

	"github.com/connapotae/hydro-shop/config"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

type IServer interface {
	Start()
}

type server struct {
	app *fiber.App
	cfg config.IConfig
	db  *sqlx.DB
}

func NewServer(cfg config.IConfig, db *sqlx.DB) IServer {
	return &server{
		cfg: cfg,
		db:  db,
		app: fiber.New(fiber.Config{
			AppName:      cfg.App().Name(),
			BodyLimit:    cfg.App().BodyLimit(),
			ReadTimeout:  cfg.App().ReadTimeout(),
			WriteTimeout: cfg.App().WriteTimeout(),
			JSONEncoder:  json.Marshal,
			JSONDecoder:  json.Unmarshal,
		}),
	}
}

func (s *server) Start() {
	//Middlewares
	middleware := InitMiddlewares(s)
	s.app.Use(middleware.Logger())
	s.app.Use(middleware.Cors())

	//Modules
	v1 := s.app.Group("api/v1")
	modules := InitModule(v1, s, middleware)

	modules.MonitorModule()
	modules.UsersModule()
	modules.AppinfoModule()

	s.app.Use(middleware.RouterCheck())

	//Graceful Shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		log.Println("server is shutting down...")
		_ = s.app.Shutdown()
	}()

	// Listen to host:port
	log.Printf("server is starting on %v", s.cfg.App().Url())
	s.app.Listen(s.cfg.App().Url())
}
