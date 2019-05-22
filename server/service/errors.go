package service

import (
	"fmt"
	"github.com/labstack/echo"
	"net/http"
)

func NewKeyNotFoundError(key string) error {
	return echo.NewHTTPError(http.StatusNotFound, fmt.Sprintf("the %s key could not be found", key))
}
