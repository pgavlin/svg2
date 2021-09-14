package svg

import (
	"encoding/xml"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func parseTestdata(path string) (*SVG, error) {
	f, err := os.Open(filepath.Join("testdata", path))
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var svg SVG
	if err = xml.NewDecoder(f).Decode(&svg); err != nil {
		return nil, err
	}
	return &svg, nil
}

func TestBadge(t *testing.T) {
	_, err := parseTestdata("badge.svg")
	assert.NoError(t, err)
}
