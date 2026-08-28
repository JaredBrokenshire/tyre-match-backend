package service_tests

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"tyre-match-backend/services"
)

func TestFileService_FileExtensionAllowed(t *testing.T) {
	s := services.NewFileService(nil)

	cases := []struct {
		Name       string
		FileName   string
		Extensions []string
		Expected   bool
	}{
		{
			Name:       "Allows matching extension",
			FileName:   "image.png",
			Extensions: []string{".png", ".jpg"},
			Expected:   true,
		},
		{
			Name:       "Rejects non matching extension",
			FileName:   "image.pdf",
			Extensions: []string{".png", ".jpg"},
			Expected:   false,
		},
		{
			Name:       "Handles empty extension list",
			FileName:   "image.png",
			Extensions: []string{},
			Expected:   false,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t,
				test.Expected,
				s.FileExtensionAllowed(test.FileName, test.Extensions),
			)
		})
	}
}

func TestFileService_MIMETypeAllowed(t *testing.T) {
	s := services.NewFileService(nil)

	cases := []struct {
		Name     string
		MimeType string
		Allowed  []string
		Expected bool
	}{
		{
			Name:     "Allows valid mime type",
			MimeType: "image/png",
			Allowed: []string{
				"image/png",
				"image/jpeg",
			},
			Expected: true,
		},
		{
			Name:     "Rejects invalid mime type",
			MimeType: "application/pdf",
			Allowed: []string{
				"image/png",
				"image/jpeg",
			},
			Expected: false,
		},
		{
			Name:     "Handles empty allowed list",
			MimeType: "image/png",
			Allowed:  []string{},
			Expected: false,
		},
	}

	for _, test := range cases {
		t.Run(test.Name, func(t *testing.T) {
			assert.Equal(t,
				test.Expected,
				s.MIMETypeAllowed(test.MimeType, test.Allowed),
			)
		})
	}
}
