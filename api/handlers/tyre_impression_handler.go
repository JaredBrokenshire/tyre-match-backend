package handlers

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"net/http"
	"strconv"
	"tyre-match-backend/api"
	"tyre-match-backend/api/requests"
	"tyre-match-backend/api/responses"
	"tyre-match-backend/services"
)

type TyreImpressionHandler struct {
	server *api.Server
}

func NewTyreImpressionHandler(server *api.Server) *TyreImpressionHandler {
	return &TyreImpressionHandler{server: server}
}

// List godoc
// @Summary List tyre impressions
// @Description List tyre impressions (paginated)
// @ID tyre-impressions-list
// @Tags TyreImpression Actions
// @Accept json
// @Produce json
// @Param page query int false "The page number"
// @Param page_size query int false "The numbers of items to return. Max 100"
// @Success 200 {object} responses.TyreImpressionPaginatedResponse
// @Router /tyre-impressions [get]
func (h *TyreImpressionHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	size, _ := strconv.Atoi(c.QueryParam("page_size"))

	tyreImpressions, count, page, size := h.server.Services.TyreImpression.List(page, size)

	res := responses.NewTyreImpressionPaginatedResponse(tyreImpressions, count, page, size)
	return responses.Response(c, http.StatusOK, res)
}

// Get godoc
// @Summary Get tyre impression by ID
// @Description Get tyre impression by ID
// @ID tyre-impressions-get
// @Tags TyreImpression Actions
// @Accept json
// @Produce json
// @Param id path int true "TyreImpression ID"
// @Success 200 {object} responses.TyreImpressionResponse
// @Failure 404 {object} responses.Error
// @Router /tyre-impressions/{id} [get]
func (h *TyreImpressionHandler) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
	}

	tyreImpression, err := h.server.Services.TyreImpression.Get(uint(id))
	switch {
	case errors.Is(err, services.NotFoundError):
		return responses.ErrorResponse(c, http.StatusNotFound, "TyreImpression not found")
	case err != nil:
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error getting TyreImpression")
	}

	res := responses.NewTyreImpressionResponse(tyreImpression)
	return responses.Response(c, http.StatusOK, res)
}

// Create godoc
// @Summary Create tyre impression
// @Description Create tyre impression
// @ID tyre-impressions-create
// @Tags TyreImpression Actions
// @Accept json
// @Produce json
// @Param params body requests.CreateTyreImpressionRequest true "Tyre Impression information"
// @Success 201 {object} responses.TyreImpressionResponse
// @Failure 400 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /tyre-impressions [post]
func (h *TyreImpressionHandler) Create(c echo.Context) error {
	req := new(requests.CreateTyreImpressionRequest)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Required fields are empty or not valid: %v", err))
	}

	dto := services.TyreImpressionDTO{
		PixelsPerInch: req.PixelsPerInch,
	}

	tyreImpression, err := h.server.Services.TyreImpression.Create(dto)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error creating TyreImpression")
	}

	res := responses.NewTyreImpressionResponse(tyreImpression)
	return responses.Response(c, http.StatusCreated, res)
}

// Upload godoc
// @Summary Upload TyreImpression image
// @Description Upload TyreImpression image
// @ID tyre-impressions-upload
// @Tags TyreImpression File Actions
// @Accept json
// @Produce json
// @Success 200 {object} responses.Data
// @Failure 400 {object} responses.Error
// @Failure 404 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /tyre-impressions/:id [post]
func (h *TyreImpressionHandler) Upload(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
	}

	file, err := requests.MultipartFileRequest(c.Request(), "file")
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Error parsing file from request")
	}

	dto := services.UploadedFile{
		Name:        file.Name,
		ContentType: file.ContentType,
		Size:        file.Size,
		Open:        file.Open,
	}
	tyreImpression, err := h.server.Services.TyreImpression.Upload(uint(id), dto)
	switch {
	case errors.Is(err, services.NotFoundError):
		log.Errorf("tyre impression not found: %v", err)
		return responses.ErrorResponse(c, http.StatusNotFound, "TyreImpression not found")
	case errors.Is(err, services.AlreadyExistsError):
		log.Errorf("tyre impression image already exists: %v", err)
		return responses.ErrorResponse(c, http.StatusBadRequest, "TyreImpression image already exists")
	case errors.Is(err, services.InvalidUploadError):
		log.Errorf("invalid image upload for tyre impression: %v", err)
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid image upload for TyreImpression")
	case errors.Is(err, services.FileStoreError):
		log.Errorf("error saving image to file store: %v", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error saving image to file store")
	case errors.Is(err, services.ProcessingError):
		log.Errorf("error processing tyre impression image: %v", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error processing TyreImpression image")
	case err != nil:
		log.Errorf("error uploading tyre impression image: %v", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error uploading TyreImpression image")
	}

	res := responses.NewTyreImpressionResponse(tyreImpression)
	return responses.Response(c, http.StatusOK, res)
}

// Delete godoc
// @Summary Delete tyre impression
// @Description Delete tyre impression
// @ID tyre-impressions-delete
// @Tags TyreImpression Actions
// @Accept json
// @Produce json
// @Param id path string true "TyreImpression ID"
// @Success 200 {object} responses.Data
// @Failure 400 {object} responses.Error
// @Failure 404 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /tyre-impressions/{id} [delete]
func (h *TyreImpressionHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
	}

	err = h.server.Services.TyreImpression.Delete(uint(id))
	switch {
	case errors.Is(err, services.NotFoundError):
		return responses.ErrorResponse(c, http.StatusNotFound, "TyreImpression not found")
	case err != nil:
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error deleting TyreImpression")
	}

	return responses.MessageResponse(c, http.StatusOK, "TyreImpression deleted successfully")
}
