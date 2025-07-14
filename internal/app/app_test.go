package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestUpdateConfig
func TestUpdateConfig(t *testing.T) {
	tests := []struct {
		target ServerConfigStruct
		source ServerConfigStruct
		want   ServerConfigStruct
	}{
		{
			target: ServerConfigStruct{ServerAddress: ""},
			source: ServerConfigStruct{ServerAddress: "source"},
			want:   ServerConfigStruct{ServerAddress: "source"},
		},
		{
			target: ServerConfigStruct{ServerAddress: "target"},
			source: ServerConfigStruct{ServerAddress: "source"},
			want:   ServerConfigStruct{ServerAddress: "target"},
		},
		{
			target: ServerConfigStruct{ServerAddress: "target"},
			source: ServerConfigStruct{ServerAddress: ""},
			want:   ServerConfigStruct{ServerAddress: "target"},
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			result := updateConfig(&test.target, &test.source)
			assert.Equal(t, test.want, *result)
		},
		)
	}
}
