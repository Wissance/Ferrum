package encoding

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HashPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		pwd := "qwerty"
		salt := "salt"
		encoder := NewPasswordJsonEncoder(salt)

		// Act
		hashedPwd := encoder.HashPassword(pwd)
		isMatch := encoder.IsPasswordsMatch(pwd, hashedPwd)

		// Assert
		assert.True(t, isMatch)
	})
}
