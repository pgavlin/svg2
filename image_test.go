package svg

import (
	"image"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func decodeTestdata(path string) (image.Image, string, error) {
	f, err := os.Open(filepath.Join("testdata", path))
	if err != nil {
		return nil, "", err
	}
	defer f.Close()

	return image.Decode(f)
}

func TestBadgeImage(t *testing.T) {
	_, format, err := decodeTestdata("badge.svg")
	require.NoError(t, err)
	assert.Equal(t, "svg", format)
}
