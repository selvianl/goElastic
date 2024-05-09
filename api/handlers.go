package api

import (
	"insider/api/docs"
	"net/http"

	"github.com/labstack/echo/v4"
	echoSwagger "github.com/swaggo/echo-swagger"
)

func (a *API) setupHandlers() {
	a.ec.GET("/", a.serveOK)
	a.ec.GET("/status", a.serveOK)
	a.ec.GET("/docs/*", echoSwagger.WrapHandler)

}

func (a *API) serveOK(c echo.Context) error {
	return c.String(http.StatusOK, "OK")
}

func (a *API) setupSwagger() {
	docs.SwaggerInfo.Title = "Insider API"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.Host = "localhost:8092"

	requireJSON := a.ensureContentType("application/json")
	a.ec.POST("/items/", a.listItems, requireJSON)

	{
		conf := a.ec.Group("/configs")
		conf.GET("/:id", a.getConfig)
		conf.GET("", a.listConfigs)
		conf.POST("", a.createConfig, requireJSON)
		conf.PUT("/:id", a.updateConfig)
		conf.DELETE("/:id", a.deleteConfig)
	}
}
