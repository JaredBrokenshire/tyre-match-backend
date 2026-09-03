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

type TyreModelHandler struct {
	server *api.Server
}

func NewTyreModelHandler(server *api.Server) *TyreModelHandler {
	return &TyreModelHandler{server: server}
}

// List godoc
// @Summary List tyre models
// @Description List tyre models (paginated)
// @ID tyre-models-list
// @Tags TyreModel Actions
// @Accept json
// @Produce json
// @Param search query string false "Search tyre models by name"
// @Param page query int false "The page number"
// @Param page_size query int false "The numbers of items to return. Max 100"
// @Success 200 {object} responses.TyreModelPaginatedResponse
// @Router /tyre-models [get]
func (h *TyreModelHandler) List(c echo.Context) error {
	page, _ := strconv.Atoi(c.QueryParam("page"))
	pageSize, _ := strconv.Atoi(c.QueryParam("page_size"))
	search := c.QueryParam("search")

	tyreModels, count, page, size := h.server.Services.TyreModel.List(page, pageSize, search)

	res := responses.NewTyreModelPaginatedResponse(tyreModels, count, page, size)
	return responses.Response(c, http.StatusOK, res)
}

// Get godoc
// @Summary Get tyre model by ID
// @Description Get tyre model by ID
// @ID tyre-models-get
// @Tags TyreModel Actions
// @Accept json
// @Produce json
// @Param id path int true "TyreModel ID"
// @Success 200 {object} responses.TyreModelResponse
// @Failure 404 {object} responses.Error
// @Router /tyre-models/{id} [get]
func (h *TyreModelHandler) Get(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
	}

	tyreModel, err := h.server.Services.TyreModel.Get(uint(id))
	switch {
	case errors.Is(err, services.NotFoundError):
		return responses.ErrorResponse(c, http.StatusNotFound, "TyreModel not found")
	case err != nil:
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error getting TyreModel")
	}

	res := responses.NewTyreModelResponse(tyreModel)
	return responses.Response(c, http.StatusOK, res)
}

// Create godoc
// @Summary Create tyre model
// @Description Create tyre model
// @ID tyre-models-create
// @Tags TyreModel Actions
// @Accept json
// @Produce json
// @Param params body requests.CreateTyreModelRequest true "Tyre Model information"
// @Success 201 {object} responses.TyreModelResponse
// @Failure 400 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /tyre-models [post]
func (h *TyreModelHandler) Create(c echo.Context) error {
	req := new(requests.CreateTyreModelRequest)
	if err := c.Bind(req); err != nil {
		return responses.Response(c, http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, fmt.Sprintf("Required fields are empty or not valid: %v", err))
	}

	dto := services.TyreModelDTO{
		Manufacturer:      req.Manufacturer,
		ModelName:         req.ModelName,
		WidthMm:           req.WidthMm,
		AspectRatio:       req.AspectRatio,
		RimDiameterInches: req.RimDiameterInches,
		GrooveCount:       req.GrooveCount,
		PixelsPerInch:     req.PixelsPerInch,
		ROITop:            req.ROITop,
		ROIBottom:         req.ROIBottom,
		ROILeft:           req.ROILeft,
		ROIRight:          req.ROIRight,
	}

	tyreModel, err := h.server.Services.TyreModel.Create(dto)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error creating TyreModel")
	}

	res := responses.NewTyreModelResponse(tyreModel)
	return responses.Response(c, http.StatusCreated, res)
}

// Update godoc
// @Summary Update tyre model
// @Description Update tyre model
// @ID tyre-models-update
// @Tags TyreModel Actions
// @Accept json
// @Produce json
// @Param id path string true "TyreModel ID"
// @Param params body requests.UpdateTyreModelRequest true "TyreModel information"
// @Success 200 {object} responses.TyreModelResponse
// @Failure 400 {object} responses.Error
// @Failure 404 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /tyre-models/{id} [put]
func (h *TyreModelHandler) Update(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
	}

	req := new(requests.UpdateTyreModelRequest)
	if err := c.Bind(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	if err := c.Validate(req); err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Required fields are empty or not valid: "+err.Error())
	}

	dto := services.TyreModelDTO{
		Manufacturer:      req.Manufacturer,
		ModelName:         req.ModelName,
		WidthMm:           req.WidthMm,
		AspectRatio:       req.AspectRatio,
		RimDiameterInches: req.RimDiameterInches,
		GrooveCount:       req.GrooveCount,
	}

	tyreModel, err := h.server.Services.TyreModel.Update(uint(id), dto)
	switch {
	case errors.Is(err, services.NotFoundError):
		return responses.ErrorResponse(c, http.StatusNotFound, "TyreModel not found")
	case err != nil:
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error updating TyreModel")
	}

	res := responses.NewTyreModelResponse(tyreModel)
	return responses.Response(c, http.StatusOK, res)
}

// Upload godoc
// @Summary Upload TyreModel image
// @Description Upload TyreModel image
// @ID tyre-models-upload
// @Tags TyreModel File Actions
// @Accept json
// @Produce json
// @Param id path string true "TyreModel ID"
// @Success 200 {object} responses.Data
// @Failure 400 {object} responses.Error
// @Failure 404 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /tyre-models/{id} [post]
func (h *TyreModelHandler) Upload(c echo.Context) error {
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
	tyreModel, err := h.server.Services.TyreModel.Upload(uint(id), dto)
	switch {
	case errors.Is(err, services.NotFoundError):
		log.Errorf("tyre model not found: %v", err)
		return responses.ErrorResponse(c, http.StatusNotFound, "TyreModel not found")
	case errors.Is(err, services.AlreadyExistsError):
		log.Errorf("tyre model image already exists: %v", err)
		return responses.ErrorResponse(c, http.StatusBadRequest, "TyreModel image already exists")
	case errors.Is(err, services.InvalidUploadError):
		log.Errorf("invalid image upload for tyre model: %v", err)
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid image upload for TyreModel")
	case errors.Is(err, services.FileStoreError):
		log.Errorf("error saving image to file store: %v", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error saving image to file store")
	case errors.Is(err, services.ProcessingError):
		log.Errorf("error processing tyre model image: %v", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error processing TyreModel image")
	case err != nil:
		log.Errorf("error uploading tyre model image: %v", err)
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error uploading TyreModel image")
	}

	res := responses.NewTyreModelResponse(tyreModel)
	return responses.Response(c, http.StatusOK, res)
}

// Delete godoc
// @Summary Delete tyre model
// @Description Delete tyre model
// @ID tyre-models-delete
// @Tags TyreModel Actions
// @Accept json
// @Produce json
// @Param id path string true "TyreModel ID"
// @Success 200 {object} responses.Data
// @Failure 400 {object} responses.Error
// @Failure 404 {object} responses.Error
// @Failure 500 {object} responses.Error
// @Router /tyre-models/{id} [delete]
func (h *TyreModelHandler) Delete(c echo.Context) error {
	id, err := strconv.ParseUint(c.Param("id"), 10, strconv.IntSize)
	if err != nil {
		return responses.ErrorResponse(c, http.StatusBadRequest, "Invalid ID")
	}

	err = h.server.Services.TyreModel.Delete(uint(id))
	switch {
	case errors.Is(err, services.NotFoundError):
		return responses.ErrorResponse(c, http.StatusNotFound, "TyreModel not found")
	case err != nil:
		return responses.ErrorResponse(c, http.StatusInternalServerError, "Error deleting TyreModel")
	}

	return responses.MessageResponse(c, http.StatusOK, "TyreModel deleted successfully")
}
