package api

import (
	"context"
	"encoding/json"
	"fmt"
	"insider/elastic"
	"insider/models"
	"net/http"

	"github.com/labstack/echo/v4"
)

// listItems godoc
//
//	@Summary	List Items
//	@Tags		items
//	@Param		esSearch	body		models.FilterParams	true	"search input"
//	@Param		limit		query		string				false	"pagination limit"
//	@Param		page		query		string				false	"active page"
//	@Success	200			{object}	api.PaginatedResponse{Results=[]models.ItemOutput, count=int}
//	@Router		/items/ [post]
func (a *API) listItems(c echo.Context) error {
	var filterParams models.FilterParams
	if err := c.Bind(&filterParams); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error parsing filter parameters")
	}

	// Check for field_name and operation are valid ones
	for _, condition := range filterParams.Conditions {
		err := condition.CheckFields()
		if err != nil {
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}

	page, limit, err := getPageAndSize(c, a.cfg.DefaultPageSize)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	var sortOption, sortOrder string

	activeConfig, err := a.db.GetActiveConfig()
	if err != nil {
		return err
	}

	if activeConfig != nil {
		sortOption = activeConfig.SortOption
		sortOrder = activeConfig.SortOrder
	}

	// Create query
	query, err := elastic.ConstructQuery(
		filterParams, sortOption, sortOrder, page, limit,
	)

	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	// Convert query to byte slice
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error parsing filter parameters")
	}

	searchRes := a.es.DoSearch(context.Background(), queryBytes)

	// Print the search results
	fmt.Println("Total hits:", searchRes.Hits.Total.Value)
	fmt.Println("Search results:")

	var items []models.ItemOutput
	for _, hit := range searchRes.Hits.Hits {
		item := models.ItemOutput{
			ItemID:   hit.Source["item_id"].(string),
			Name:     hit.Source["name"].(string),
			Locale:   hit.Source["locale"].(string),
			Click:    int(hit.Source["click"].(float64)),
			Purchase: int(hit.Source["purchase"].(float64)),
		}
		items = append(items, item)
	}

	return c.JSON(http.StatusOK, PaginatedResponse{
		Results: items,
		Count:   uint64(searchRes.Hits.Total.Value),
	})
}
