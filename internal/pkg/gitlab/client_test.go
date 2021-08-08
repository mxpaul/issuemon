package gitlab

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGitlabClient(t *testing.T) {
	glToken := "12345"
	client, err := NewGitlabClient(glToken)
	assert.NoError(t, err)
	assert.IsType(t, &GitlabClient{}, client)
}
