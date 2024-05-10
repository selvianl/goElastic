package models

type OperationEnum string
type FieldEnum string

const (
	Lt     OperationEnum = "lt"
	Gt     OperationEnum = "gt"
	Equals OperationEnum = "equals"
	Query  OperationEnum = "query"

	Name     FieldEnum = "name"
	Locale   FieldEnum = "locale"
	Click    FieldEnum = "click"
	Purchase FieldEnum = "purchase"
)

func GetFieldEnumValue(field FieldEnum) string {
	switch field {
	case Name:
		return "name"
	case Locale:
		return "locale"
	case Click:
		return "click"
	case Purchase:
		return "purchase"
	default:
		return ""
	}
}
