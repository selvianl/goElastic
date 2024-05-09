package models

type FilterParams struct {
	Conditions []FilterCondition `json:"conditions"`
}

type FilterCondition struct {
	FieldName string `json:"field_name"`
	Operation string `json:"operation"`
	Value     string `json:"value"`
}

type ConfigInput struct {
	SortOption string `json:"sort_option" validate:"required" err:"sort_option is required (click, purchase, name, item_id or locale)"`
	SortOrder  string `json:"sort_order" validate:"required" err:"sort_order is required (asc or desc)"`
	IsActive   bool   `json:"is_active" validate:"omitempty" err:"is_active is required (true or false)"`
}

type ConfigInputUpdate struct {
	SortOption string `json:"sort_option" validate:"omitempty"`
	SortOrder  string `json:"sort_order" validate:"omitempty"`
	IsActive   bool   `json:"is_active" validate:"omitempty"`
}
