package jsontools

import (
	"encoding/json"
	"strings"
)

func MergeNonIntersect[T1 any, T2 any](first *T1, second *T2) (interface{}, string) {
	var result interface{}
	str1, err := json.Marshal(&first)
	if err != nil {
		// todo(UMV): add trouble logging
		return nil, ""
	}
	str2, err := json.Marshal(&second)
	if err != nil {
		// todo(UMV): add trouble logging
		return nil, ""
	}
	// trim } from end of str1
	str1 = []byte(strings.TrimRight(string(str1), "}"))

	// trim { from start of str2
	str2 = []byte(strings.TrimLeft(string(str2), "{"))
	str := string(str1) + "," + string(str2)

	err = json.Unmarshal([]byte(str), &result)
	if err != nil {
		// todo(UMV): add trouble logging
		return nil, ""
	}
	return result, str
}
