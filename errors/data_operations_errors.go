package errors

import (
	"errors"
	sf "github.com/wissance/stringFormatter"
)

var (
	EmptyNotFoundErr           = ObjectNotFoundError{}
	ErrZeroLength              = errors.New("zero length")
	ErrNotAll                  = errors.New("not all values")
	ErrExists                  = ObjectAlreadyExistsError{}
	ErrNotExists               = errors.New("not exists")
	ErrOperationNotSupported   = errors.New("manager operation is not supported yet (temporarily or permanent)")
	ErrOperationNotImplemented = errors.New("manager operation is not implemented yet (wait for future releases)")
	ErrDataSourceNotAvailable  = DataProviderNotAvailable{}
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
	operation   string
	method      string
	internalErr error
}

type DataProviderNotAvailable struct {
	providerType string
	source       string
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

func NewUnknownError(operation string, method string, internalErr error) UnknownError {
	return UnknownError{operation: operation, method: method, internalErr: internalErr}
}

func (e UnknownError) Error() string {
	return sf.Format("An error occurred during: \"{0}\" in method: \"{1}\", internal error: {2}", e.operation, e.method, e.internalErr)
}

func NewDataProviderNotAvailable(providerType string, source string) DataProviderNotAvailable {
	return DataProviderNotAvailable{providerType: providerType, source: source}
}

func (e DataProviderNotAvailable) Error() string {
	return sf.Format("{0} is not ready/up/available, please try again later", e.providerType)
}
