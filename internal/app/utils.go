// package app
package app

import (
	"hash/fnv"
	"strconv"

	"github.com/hashicorp/go-uuid"
)

// func hash
func Hash(s string) string {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return ""
	}
	return strconv.Itoa(int(h.Sum32()))
}

// func GetNewUUID
func GetNewUUID() (string, error) {
	userID, err := uuid.GenerateUUID()
	return string(userID), err
}
