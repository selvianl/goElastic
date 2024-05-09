package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func (a *API) ensureContentType(contentType string) func(next echo.HandlerFunc) echo.HandlerFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ct := strings.Split(c.Request().Header.Get(echo.HeaderContentType), ";")[0]
			if ct != contentType {
				return echo.NewHTTPError(
					http.StatusBadRequest,
					fmt.Errorf("invalid content-type %q, expecting %q", ct, contentType),
				)
			}
			return next(c)
		}
	}
}
