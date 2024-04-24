package errors

import (
	"errors"
	sf "github.com/wissance/stringFormatter"
)

var (
	EmptyNotFoundErr         = ObjectNotFoundError{}
	ErrZeroLength            = errors.New("zero length")
	ErrNotAll                = errors.New("not all values")
	ErrExists                = ObjectAlreadyExistsError{}
	ErrNotExists             = errors.New("not exists")
	ErrOperationNotSupported = errors.New("manager operation is not supported yet (temporarily or permanent)")
)

type ObjectAlreadyExistsError struct {
	objectType     string
	objectId       string
	additionalInfo string
}

type ObjectNotFoundError struct {
	objectType     string
	objectId       string
	additionalInfo string
}

type UnknownError struct {
	method      string
	internalErr error
}

func NewObjectExistsError(objectType string, objectId string, additional string) ObjectAlreadyExistsError {
	return ObjectAlreadyExistsError{objectId: objectId, objectType: objectType, additionalInfo: additional}
}

func (e ObjectAlreadyExistsError) Error() string {
	return sf.Format("object of type \"{0}\" with id: \"{1}\" already exists in data store, additional data: {2}", e.objectType, e.objectId,
		e.additionalInfo)
}

func NewObjectNotFoundError(objectType string, objectId string, additional string) ObjectNotFoundError {
	return ObjectNotFoundError{objectId: objectId, objectType: objectType, additionalInfo: additional}
}

func (e ObjectNotFoundError) Error() string {
	return sf.Format("object of type \"{0}\" with id: \"{1}\" was not found in data store, additional data: {2}", e.objectType, e.objectId,
		e.additionalInfo)
}

func NewUnknownError(method string, internalErr error) UnknownError {
	return UnknownError{method: method, internalErr: internalErr}
}

func (e UnknownError) Error() string {
	return sf.Format("An error occurred in method: \"{0}\", internal error: {1}", e.method, e.internalErr)
}
