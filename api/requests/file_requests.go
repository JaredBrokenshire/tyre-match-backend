package requests

import (
	"fmt"
	"io"
	"net/http"
)

type UploadedFile struct {
	Name        string
	ContentType string
	Size        int64
	Open        func() (io.ReadCloser, error)
}

func MultipartFileRequest(r *http.Request, field string) (UploadedFile, error) {
	if err := r.ParseMultipartForm(10 * 1024 * 1024); err != nil {
		return UploadedFile{}, fmt.Errorf("error parsing multipart form file: %v", err)
	}

	f, h, err := r.FormFile(field)
	if err != nil {
		return UploadedFile{}, fmt.Errorf("error getting multipart form file:: %w", err)
	}

	f.Close()

	return UploadedFile{Name: h.Filename, Size: h.Size, Open: func() (io.ReadCloser, error) { file, _, err := r.FormFile(field); return file, err }}, nil
}
