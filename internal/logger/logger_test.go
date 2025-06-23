package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPostJSON(t *testing.T) {
	err := InitLogger("DEBUG")
	assert.Nil(t, err)
	err = InitLogger("INFO")
	assert.Nil(t, err)

}
