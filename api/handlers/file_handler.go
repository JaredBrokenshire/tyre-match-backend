package handlers

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"tyre-match-backend/api"
	"tyre-match-backend/api/responses"
)

type FileHandler struct {
	server *api.Server
}

func NewFileHandler(server *api.Server) *FileHandler {
	return &FileHandler{server: server}
}

// Get godoc
// @Summary Get file
// @Description Get file from file store
// @ID files-get
// @Tags Files Actions
// @Produce text/plain
// @Param filepath path string true "filepath"
// @Success 200 {file} runtime.File
// @Failure 404 {object} responses.Error
// @Security ApiKeyAuth
// @Router /files/{filepath} [get]
func (h *FileHandler) Get(c echo.Context) error {
	path := c.Param("*")
	data, err := h.server.Dependencies.GetFileStore().ReadFile(path)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Unable to read file")
	}

	return c.Blob(http.StatusOK, http.DetectContentType(data), data)
}
