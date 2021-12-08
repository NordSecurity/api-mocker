package main

import (
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"time"
)

const (
	green   = "\033[97;42m"
	white   = "\033[90;47m"
	yellow  = "\033[90;43m"
	red     = "\033[97;41m"
	blue    = "\033[97;44m"
	magenta = "\033[97;45m"
	cyan    = "\033[97;46m"
	reset   = "\033[0m"
)

func FormattedLog(c echo.Context) {
	fmt.Printf("%v |%s %3d %s| %15v | %s %-7s %s %#v\n",
		time.Now().Format("2006/01/02 - 15:04:05"),
		statusCodeColor(c.Response().Status), c.Response().Status, reset,
		c.RealIP(),
		methodColor(c.Request().Method), c.Request().Method, reset,
		c.Request().RequestURI,
	)
}

// statusCodeColor is the ANSI color for appropriately logging http status code to a terminal.
func statusCodeColor(code int) string {
	switch {
	case code >= http.StatusOK && code < http.StatusMultipleChoices:
		return green
	case code >= http.StatusMultipleChoices && code < http.StatusBadRequest:
		return white
	case code >= http.StatusBadRequest && code < http.StatusInternalServerError:
		return yellow
	default:
		return red
	}
}

// methodColor is the ANSI color for appropriately logging http method to a terminal.
func methodColor(method string) string {
	switch method {
	case http.MethodGet:
		return blue
	case http.MethodPost:
		return cyan
	case http.MethodPut:
		return yellow
	case http.MethodDelete:
		return red
	case http.MethodPatch:
		return green
	case http.MethodHead:
		return magenta
	case http.MethodOptions:
		return white
	default:
		return reset
	}
}
