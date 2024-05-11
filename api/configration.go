package api

import (
	"errors"
	"insider/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// createConfig godoc
//
//	@Summary	Create config
//	@Tags		configs
//	@Accept		json
//	@Produce	json
//	@Param		applicant	body		models.ConfigInput	true	"Create Config"
//	@Success	201			{object}	models.SortConfig
//	@Router		/configs [post]
func (a *API) createConfig(c echo.Context) error {
	var in *models.ConfigInput
	if err := c.Bind(&in); err != nil {
		return err
	}

	if err := a.val.Struct(in); err != nil {
		return err
	}

	if in.IsActive {
		activeConfig, err := a.db.GetActiveConfig()
		if err != nil {
			return err
		}

		if activeConfig != nil {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("there is already active conf exists"))
		}
	}

	conf, err := a.db.CreateConfig(*in)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusCreated, conf)
}

// listConfig godoc
//
//	@Summary	List configs
//	@Tags		configs
//	@Accept		json
//	@Produce	json
//	@Param		limit	query		string	false	"pagination limit"
//	@Param		page	query		string	false	"active page"
//	@Success	200		{object}	PaginatedResponse{Results=[]models.SortConfig, count=int}
//	@Router		/configs [get]
func (a *API) listConfigs(c echo.Context) error {
	offset, limit, err := getPageAndSize(c, 20)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	count, results, err := a.db.ListConfigs(int(offset), int(limit))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, PaginatedResponse{
		Results: &results,
		Count:   uint64(count),
	})
}

// getConfig godoc
//
//	@Summary	Get config
//	@Tags		configs
//	@Accept		json
//	@Produce	json
//	@Param		id	path		string	true	"Config id"
//	@Success	200	{object}	models.SortConfig
//	@Router		/configs/{id} [get]
func (a *API) getConfig(c echo.Context) error {
	sc, err := a.db.GetConfig(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if sc == nil {
		return echo.NewHTTPError(http.StatusNotFound, errors.New("config not found"))
	}

	return c.JSON(http.StatusOK, sc)
}

// updateConfig godoc
//
//	@Summary	Update config
//	@Tags		configs
//	@Accept		json
//	@Produce	json
//	@Param		id		path		string						true	"Config ID"
//	@Param		address	body		models.ConfigInputUpdate	true	"Update Config"
//	@Success	200		{object}	models.SortConfig
//	@Router		/configs/{id} [put]
func (a *API) updateConfig(c echo.Context) error {
	var in models.ConfigInputUpdate
	if err := c.Bind(&in); err != nil {
		return err
	}

	if err := a.val.Struct(in); err != nil {
		return err
	}

	if in.IsActive {
		activeConfig, err := a.db.GetActiveConfig()
		if err != nil {
			return err
		}

		if activeConfig != nil {
			return echo.NewHTTPError(http.StatusBadRequest, errors.New("there is already active conf exists"))
		}
	}

	cfg, err := a.db.GetConfig(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if cfg == nil {
		return echo.NewHTTPError(http.StatusNotFound, errors.New("config not found"))
	}

	cfg, err = a.db.UpdateConfig(cfg, in)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.JSON(http.StatusOK, cfg)
}

// deleteConfig godoc
//
//	@Summary	Delete Config
//	@Tags		configs
//	@Accept		json
//	@Param		id	path	string	true	"Config id"
//	@Success	204
//	@Router		/configs/{id} [delete]
func (a *API) deleteConfig(c echo.Context) error {
	sc, err := a.db.GetConfig(c.Param("id"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	if err := a.db.DeleteConfig(sc); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusNoContent)
}
