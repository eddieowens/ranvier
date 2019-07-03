package model

import "github.com/labstack/echo"

type Route struct {
	Method      string
	Path        string
	HandlerFunc echo.HandlerFunc
}

func NewRoute(method string, path string, handler echo.HandlerFunc) Route {
	return Route{
		Method:      method,
		Path:        path,
		HandlerFunc: handler,
	}
}
