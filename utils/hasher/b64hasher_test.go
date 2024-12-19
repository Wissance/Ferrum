package b64hasher

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_HashPassword(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		// Arrange
		pwd := "qwerty"
		salt := "salt"

		// Act
		hashedPwd := HashPassword(pwd, salt)
		isMatch := IsPasswordsMatch(pwd, salt, hashedPwd)

		// Assert
		assert.True(t, isMatch)
	})
}
