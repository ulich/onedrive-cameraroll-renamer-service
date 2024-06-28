package internal

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCalcFilename(t *testing.T) {
	for _, name := range []string{
		"20190519_205940.jpg",
		"20190519_205940.jpeg",
		"20190519_205940.mp4",
		"20190519_205940_11212.jpg",
		"20190519_205940_11212.jpeg",
		"20190519_205940_11212.mp4",
		"20190519_205940(0).jpg",
		"20190519_205940(0).jpeg",
		"20190519_205940(0).mp4",
		"20190519_205940(2).jpg",
	} {
		t.Run("keeps a filename the same if its already in the desired format: "+name, func(t *testing.T) {
			newFilename := doCalcNewFilename(t, name)
			assert.Equal(t, name, newFilename)
		})
	}

	for _, tt := range []testCase{
		{"20190519_205940.JPG", "20190519_205940.jpg"},
		{"20190519_205940.JPEG", "20190519_205940.jpeg"},
		{"20190519_205940.MP4", "20190519_205940.mp4"},
		{"20190519_205940.JPG", "20190519_205940.jpg"},
		{"20190519_205940.JPEG", "20190519_205940.jpeg"},
		{"20190519_205940.MP4", "20190519_205940.mp4"},
	} {
		t.Run("renames a filename into lower case: "+tt.name, func(t *testing.T) {
			newFilename := doCalcNewFilename(t, tt.name)
			assert.Equal(t, tt.expectedName, newFilename)
		})
	}

	for _, tt := range []testCase{
		{"IMG_20190519_205940.jpg", "20190519_205940.jpg"},
		{"IMG_20190519_205940.jpeg", "20190519_205940.jpeg"},
	} {
		t.Run("renames a filename with IMG_ prefix: "+tt.name, func(t *testing.T) {
			newFilename := doCalcNewFilename(t, tt.name)
			assert.Equal(t, tt.expectedName, newFilename)
		})
	}

	for _, tt := range []testCase{
		{"IMG_1558292380001_12345.jpg", "20190519_205940_12345.jpg"},
		{"IMG_1558292380001_12345.jpeg", "20190519_205940_12345.jpeg"},
	} {
		t.Run("renames a filename with IMG_ prefix and timestamp: "+tt.name, func(t *testing.T) {
			newFilename := doCalcNewFilename(t, tt.name)
			assert.Equal(t, tt.expectedName, newFilename)
		})
	}

	for _, tt := range []testCase{
		{"IMG-20190519-WA1234.jpg", "20190519_1234_WA.jpg"},
		{"IMG-20190519-WA1234.jpeg", "20190519_1234_WA.jpeg"},
	} {
		t.Run("renames a filename with IMG_ prefix and -WA suffix: "+tt.name, func(t *testing.T) {
			newFilename := doCalcNewFilename(t, tt.name)
			assert.Equal(t, tt.expectedName, newFilename)
		})
	}

	t.Run("renames a filename with VID_ prefix", func(t *testing.T) {
		newFilename := doCalcNewFilename(t, "VID_20190519_205940.mp4")
		assert.Equal(t, "20190519_205940.mp4", newFilename)
	})

	t.Run("renames a filename with VID- prefix and -WA suffix", func(t *testing.T) {
		newFilename := doCalcNewFilename(t, "VID-20190519-WA1234.mp4")
		assert.Equal(t, "20190519_1234_WA.mp4", newFilename)
	})

	t.Run(`returns an error for all other filenames`, func(t *testing.T) {
		_, err := calcNewFilename("foo")
		if err == nil {
			t.Errorf("expected an error for unknown filename pattern")
		}
	})
}

type testCase struct {
	name         string
	expectedName string
}

func doCalcNewFilename(t *testing.T, filename string) string {
	newName, err := calcNewFilename(filename)
	require.NoError(t, err)
	return newName
}
