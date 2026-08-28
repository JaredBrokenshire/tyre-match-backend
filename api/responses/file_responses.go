package responses

import (
	"time"
	"tyre-match-backend/db/models"
)

type FileResponse struct {
	ID uint `json:"id"`

	Model    string `json:"model"`
	ModelId  uint   `json:"model_id"`
	FileType string `json:"file_type"`

	Name     string `json:"name"`
	Location string `json:"location"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func NewFileResponse(file *models.File) *FileResponse {
	return &FileResponse{
		ID:        file.ID,
		Model:     file.Model,
		ModelId:   file.ModelId,
		FileType:  file.FileType,
		Name:      file.Name,
		Location:  file.Location,
		CreatedAt: file.CreatedAt,
		UpdatedAt: file.UpdatedAt,
	}
}

func NewFileResponses(files []*models.File) []FileResponse {
	var res []FileResponse
	for _, file := range files {
		res = append(res, *NewFileResponse(file))
	}
	return res
}
