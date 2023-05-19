package validators

import "strconv"

type ValueTypeRequirements string

const (
	String   ValueTypeRequirements = "string"
	Integer  ValueTypeRequirements = "integer"
	Boolean  ValueTypeRequirements = "boolean"
	StrOrInt ValueTypeRequirements = "str or int"
	Any      ValueTypeRequirements = "any"
)

func IsStrValueOfRequiredType(requirements ValueTypeRequirements, value *string) bool {
	if value == nil {
		return false
	}
	if requirements == Any {
		return true
	}
	if requirements == Integer {
		_, err := strconv.Atoi(*value)
		return err == nil
	}
	if requirements == Boolean {
		_, err := strconv.ParseBool(*value)
		return err == nil
	}
	return true
}
