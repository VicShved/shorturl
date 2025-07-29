package middware

import (
	"testing"

	"github.com/VicShved/shorturl/internal/app"
	"github.com/stretchr/testify/assert"
)

func TestAuth(t *testing.T) {
	userID, err := app.GetNewUUID()
	assert.Nil(t, err)
	token, err := GetJWTTokenString(&userID)
	assert.Nil(t, err)
	_, resultUserID, err := ParseTokenUserID(token)
	assert.Nil(t, err)
	assert.Equal(t, userID, resultUserID)
}
