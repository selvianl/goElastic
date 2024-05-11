package models

import (
	"fmt"
	"strconv"
	"unicode"
)

type ConfigInput struct {
	SortOption string `json:"sort_option" validate:"required"`
	SortOrder  string `json:"sort_order" validate:"required"`
	IsActive   bool   `json:"is_active" validate:"omitempty"`
}

type ConfigInputUpdate struct {
	SortOption string `json:"sort_option" validate:"omitempty"`
	SortOrder  string `json:"sort_order" validate:"omitempty"`
	IsActive   bool   `json:"is_active" validate:"omitempty"`
}
type FilterParams struct {
	Conditions []FilterCondition `json:"conditions"`
}

type FilterCondition struct {
	FieldName FieldEnum     `json:"field_name"`
	Operation OperationEnum `json:"operation"`
	Value     string        `json:"value"`
}

func (q *FilterCondition) CheckFields() error {
	switch q.FieldName {
	case Name, Locale:
		if q.Operation == "query" {
			if isAlphabetic(q.Value) {
				return nil
			} else {
				return fmt.Errorf("`value` is not alphabetic")
			}
		} else {
			return fmt.Errorf("if field_name is `name` or `locale` operation must be `query` and `value` must be something alphabetic")
		}
	case Click, Purchase:
		if q.Operation == "lt" || q.Operation == "gt" || q.Operation == "equals" {
			if isNumeric(q.Value) {
				return nil
			} else {
				return fmt.Errorf("`value` is not digit")
			}
		} else {
			return fmt.Errorf("if field_name is `click` or `purchase` operation must `lt` or `gt` or `equals` and value must be something only digits")
		}
	default:
		return fmt.Errorf("unsupported field_name: %s", q.FieldName)
	}
}

func isAlphabetic(value string) bool {
	for _, ch := range value {
		if !unicode.IsLetter(ch) {
			return false
		}
	}
	return true
}

func isNumeric(value string) bool {
	_, err := strconv.Atoi(value)
	return err == nil
}
