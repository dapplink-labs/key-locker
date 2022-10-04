package api

import (
	"github.com/labstack/echo/v4"
	"net/http"
)

func (s *Server) GetKeyHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}

func (s *Server) SetKeyHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, "")
}
