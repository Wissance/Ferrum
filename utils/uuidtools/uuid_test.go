package uuidtools

import (
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsUUIDEmpty(t *testing.T) {
	createdUuid := uuid.New()
	for name, test := range map[string]struct {
		identifier      *uuid.UUID
		expectedIsEmpty bool
	}{
		"nil pointer uuid": {
			identifier:      nil,
			expectedIsEmpty: true,
		},
		"non initialized uuid": {
			identifier:      &uuid.UUID{},
			expectedIsEmpty: true,
		},
		"created uuid": {
			identifier:      &createdUuid,
			expectedIsEmpty: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			actualIsEmpty := IsUUIDEmpty(test.identifier)
			assert.Equal(t, test.expectedIsEmpty, actualIsEmpty)
		})
	}
}

func TestGetOrCreateUUID(t *testing.T) {
	id := GetOrCreateNewNonEmptyUuid(nil)
	assert.NotNil(t, id)
	idCopy := id
	idCopy = GetOrCreateNewNonEmptyUuid(&idCopy)
	assert.Equal(t, id, idCopy)
}
