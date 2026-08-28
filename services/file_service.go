package services

import (
	"bytes"
	"fmt"
	"path/filepath"
	"slices"
	"strings"
	"tyre-match-backend/pkg/file_storage"
)

type SaveFileRequest struct {
	Data            []byte
	Name            string
	TargetDirectory string
	Model           string
	ModelId         uint
	FileType        string
	Extension       string
}

type FileService struct {
	validator *UploadedFileValidator
	FileStore file_storage.Store
}

func NewFileService(fileStore file_storage.Store) *FileService {
	return &FileService{
		validator: NewUploadedFileValidator(),
		FileStore: fileStore,
	}
}

func (s *FileService) Validate(file UploadedFile) error {
	return s.validator.Validate(file)
}

func (s *FileService) FileExtensionAllowed(fileName string, extensions []string) bool {
	ext := strings.TrimPrefix(strings.ToLower(filepath.Ext(fileName)), ".")
	for _, e := range extensions {
		if strings.TrimPrefix(strings.ToLower(e), ".") == ext {
			return true
		}
	}
	return false
}

func (s *FileService) MIMETypeAllowed(mime string, allowed []string) bool {
	return slices.Contains(allowed, mime)
}

func (s *FileService) SaveFile(request SaveFileRequest) error {
	err := s.validateSaveFileRequest(request)
	if err != nil {
		return fmt.Errorf("error validating save file request: %w", err)
	}

	name := strings.TrimSuffix(
		request.Name,
		filepath.Ext(request.Name),
	)
	extension := strings.TrimPrefix(
		strings.ToLower(request.Extension),
		".",
	)

	err = s.FileStore.Save(
		request.TargetDirectory,
		bytes.NewReader(request.Data),
		fmt.Sprintf("%s.%s", name, extension),
	)
	if err != nil {
		return fmt.Errorf("error saving file to file store: %w", err)
	}

	return nil
}

func (s *FileService) validateSaveFileRequest(request SaveFileRequest) error {
	if request.Extension == "" {
		return fmt.Errorf("file extension is required")
	}
	if len(request.Data) == 0 {
		return fmt.Errorf("file data is empty")
	}
	if request.Name == "" {
		return fmt.Errorf("file name is required")
	}
	if request.TargetDirectory == "" {
		return fmt.Errorf("target directory is required")
	}
	if request.Model == "" {
		return fmt.Errorf("model is required")
	}
	if request.ModelId == 0 {
		return fmt.Errorf("model id is required")
	}
	if request.FileType == "" {
		return fmt.Errorf("file type is required")
	}
	return nil
}
