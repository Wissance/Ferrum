package data

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	sf "github.com/wissance/stringFormatter"
	"testing"
)

func TestInitUserWithJsonAndCheck(t *testing.T) {
	testCases := []struct {
		name              string
		userName          string
		preferredUsername string
	}{
		{name: "simple_user_data", userName: "admin", preferredUsername: "Administrator"},
	}

	for _, tCase := range testCases {
		t.Run(tCase.name, func(t *testing.T) {
			t.Parallel()
			jsonTemplate := `{"info":{"name":"{0}", "preferred_username": "{1}"}}`
			jsonStr := sf.Format(jsonTemplate, tCase.userName, tCase.preferredUsername)
			var rawUserData interface{}
			err := json.Unmarshal([]byte(jsonStr), &rawUserData)
			assert.NoError(t, err)
			user := CreateUser(rawUserData)
			assert.Equal(t, tCase.preferredUsername, user.GetUsername())
		})
	}
}
