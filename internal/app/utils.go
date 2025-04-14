package app

import (
	"hash/fnv"
	"strconv"

	"github.com/hashicorp/go-uuid"
)

// hash
func Hash(s string) string {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return ""
	}
	return strconv.Itoa(int(h.Sum32()))
}

func GetNewUUID() (TypeUserID, error) {
	userID, err := uuid.GenerateUUID()
	return TypeUserID(userID), err
}
