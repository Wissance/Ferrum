package uuidtools

import "github.com/google/uuid"

func IsUUIDEmpty(uniqueIdentifier *uuid.UUID) bool {
	if uniqueIdentifier == nil {
		return true
	}
	uuidBytes, _ := uniqueIdentifier.MarshalBinary()
	isEmpty := true
	for _, b := range uuidBytes {
		if b != 0 {
			isEmpty = false
			break
		}
	}
	return isEmpty
}

func GetOrCreateNewNonEmptyUuid(uniqueIdentifier *uuid.UUID) uuid.UUID {
	return uuid.New()
}
