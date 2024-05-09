package models

type ItemOutput struct {
	ItemID   string `json:"item_id"`
	Name     string `json:"name"`
	Locale   string `json:"locale"`
	Click    int    `json:"click"`
	Purchase int    `json:"purchase"`
}

type SearchResponse struct {
	Hits struct {
		Total struct {
			Value int `json:"value"`
		} `json:"total"`
		Hits []struct {
			Source map[string]interface{} `json:"_source"`
		} `json:"hits"`
	} `json:"hits"`
}
