package gh

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetLatestReleaseTagName(t *testing.T) {
	tagName, err := GetLatestReleaseTagName("api.github.com", "", "traefik", "structor")

	require.NoError(t, err)
	assert.Regexp(t, `v\d+.\d+(.\d+)?`, tagName)
}

func TestGetLatestReleaseTagName_Errors(t *testing.T) {
	_, err := GetLatestReleaseTagName("api.github.com", "", "error", "error")

	assert.EqualError(t, err, `failed to get latest release tag name on GitHub ("https://api.github.com/repos/error/error/releases/latest"), status: 404 Not Found`)
}
