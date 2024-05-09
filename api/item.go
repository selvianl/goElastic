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
	page, limit, err := getPageAndSize(c, 20)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	query := elastic.ConstructQuery(filterParams)
	query["from"] = page
	query["size"] = limit

	// Convert query to byte slice
	queryBytes, err := json.Marshal(query)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Error parsing filter parameters")
	}

	searchRes := elastic.DoSearch(context.Background(), queryBytes, a.es)

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
