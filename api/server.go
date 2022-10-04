package api

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/savour-labs/key-locker/config"
	"gorm.io/gorm"
	"strconv"
)

type Server struct {
	db   *gorm.DB
	echo *echo.Echo
	port int
}

func NewServer(db *gorm.DB, cfg *config.Server) *Server {
	e := echo.New()
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "method=${method}, uri=${uri}, status=${status}\n",
	}))
	e.Use(middleware.Recover())
	e.Debug = cfg.Debug
	server := &Server{
		db:   db,
		echo: e,
		port: cfg.Port,
	}
	server.routes()
	return server
}

func (s *Server) routes() {
	s.echo.GET("ket/:get", s.GetKeyHandler)
	s.echo.GET("ket/:set", s.SetKeyHandler)
}

func (s *Server) Run() {
	s.echo.Logger.Fatal(s.echo.Start(":" + strconv.Itoa(s.port)))
}
