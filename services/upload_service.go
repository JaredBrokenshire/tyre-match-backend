package services

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path/filepath"
	"slices"
	"strings"
)

type UploadedFile struct {
	Name        string
	ContentType string
	Size        int64
	Open        func() (io.ReadCloser, error)
}

type FileValidator interface {
	Validate(file UploadedFile) error
}

type UploadedFileValidator struct {
	AllowedExtensions []string
	AllowedMIMEs      []string
	MaxBytes          int64
}

func NewUploadedFileValidator() *UploadedFileValidator {
	return &UploadedFileValidator{
		AllowedExtensions: []string{"jpg", "jpeg", "png", "webp"},
		AllowedMIMEs:      []string{"image/jpeg", "image/png", "image/webp"},
	}
}

func (v *UploadedFileValidator) Validate(file UploadedFile) error {
	if file.Open == nil {
		return fmt.Errorf("%w: file cannot be opened", InvalidUploadError)
	}
	if file.Name == "" {
		return fmt.Errorf("%w: filename is empty", InvalidUploadError)
	}

	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(file.Name)), ".")
	if !slices.Contains(v.AllowedExtensions, ext) {
		return fmt.Errorf("%w: invalid file extension", InvalidUploadError)
	}

	f, err := file.Open()
	if err != nil {
		return fmt.Errorf("%w: unable to open file: %v", InvalidUploadError, err)
	}
	defer f.Close()

	header := make([]byte, 512)
	n, err := io.ReadFull(f, header)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) {
		return fmt.Errorf("%w: unable to read file: %v", InvalidUploadError, err)
	}
	if n == 0 {
		return fmt.Errorf("%w: file is empty", InvalidUploadError)
	}

	actualMIME := http.DetectContentType(header[:n])
	if !slices.Contains(v.AllowedMIMEs, actualMIME) {
		return fmt.Errorf("%w: invalid file content type", InvalidUploadError)
	}

	return nil
}
